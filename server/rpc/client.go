package rpc

import (
	"customermanager-go/common/logger"
	"customermanager-go/common/proto/rpc/auth"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"sync"
)

var (
	once       sync.Once
	AuthClient auth.AuthService
)

func InitMicro(name string, version string) {
	once.Do(func() {
		reg := consul.NewRegistry(func(options *registry.Options) {
			options.Addrs = []string{"127.0.0.1:8500"}
		})
		service := micro.NewService(micro.Name(name), micro.Version(version), micro.Registry(reg), micro.Address("127.0.0.1:61008"))
		service.Init()

		AuthClient = auth.NewAuthService("auth", service.Client())
		go func() {
			if err := service.Run(); err != nil {
				logger.Error("service auth run failed, err: %s", err.Error())
			}
		}()
	})
}
