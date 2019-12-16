package main

import (
	service2 "customermanager-go/auth/service"
	"customermanager-go/common/db"
	"customermanager-go/common/logger"
	"customermanager-go/common/proto/rpc/auth"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
)

const (
	APP_NAME = "auth"
	VERSION  = "v1"
)

func main() {
	logger.StartLogger("auth.log", "info")

	dbString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&timeout=10s&loc=Local&parseTime=true", "root", "xiaoxiong581", "192.168.137.106", 3306, "merchant")
	if err := db.InitEngine(dbString, 16*1024*1024); err != nil {
		logger.Error("fail to init db, error: %s", err.Error())
	}

	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	service := micro.NewService(micro.Name(APP_NAME), micro.Version(VERSION), micro.Registry(reg), micro.Address("127.0.0.1:60984"))
	service.Init()
	auth.RegisterAuthServiceHandler(service.Server(), new(service2.AuthService))

	if err := service.Run(); err != nil {
		logger.Error("service auth run failed, err: %s", err.Error())
	}
}
