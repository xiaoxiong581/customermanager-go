package main

import (
    "context"
    "flag"
    "fmt"
    "os"
    "os/signal"
    "time"
    "customermanager-go/server/common/db"
    "customermanager-go/server/common/logger"
    "customermanager-go/server/http"
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

    go func() {
        if err := httpServer.ListenAndServe(); err != nil {
            logger.Error("listen http server fail, error: %s", err.Error())
        }
    }()

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
