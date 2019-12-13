package main

import (
	"context"
	"customermanager-go/server/cron"
	"customermanager-go/server/db"
	"customermanager-go/server/http"
	"customermanager-go/server/logger"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
)

const (
	APIName = "customermanager-go"
)

var (
	ip   = flag.String("ip", "0.0.0.0", "server ip")
	port = flag.Int("port", 29080, "server port")
)

func main() {
	dbString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&timeout=10s&loc=Local&parseTime=true", "root", "xiaoxiong581", "192.168.137.106", 3306, "merchant")
	if err := db.InitEngine(dbString, 16*1024*1024); err != nil {
		logger.Error("fail to init db, error: %s", err.Error())
	}

	//serverCtx := utils.NewContext(utils.NewUUID(), "")
	httpServer := http.NewHttpServer(&http.ServerConfig{
		IP:   *ip,
		Port: *port,
	})

	appRootPath, _ := os.Getwd()
	logger.Info("get app root path: %s", appRootPath)
	crtPath := strings.Join([]string{appRootPath, "config", "cert", "server.crt"}, string(filepath.Separator))
	keyPath := strings.Join([]string{appRootPath, "config", "cert", "server.key"}, string(filepath.Separator))

	go func() {
		if err := httpServer.ListenAndServeTLS(crtPath, keyPath); err != nil {
			logger.Error("listen http server fail, error: %s", err.Error())
		}
	}()

	cron.StartCron()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error(err.Error())
	}
	logger.Info("Success shutting server.")
}
