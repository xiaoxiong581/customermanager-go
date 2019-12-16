package rpc

import (
	"customermanager-go/common/logger"
	"customermanager-go/common/proto/rpc/auth"
	"github.com/micro/go-micro"
	"sync"
)

var (
	once       sync.Once
	AuthClient auth.AuthService
)

func InitMicro(name string, version string) {
	once.Do(func() {
		service := micro.NewService(micro.Name(name), micro.Version(version), micro.Address("127.0.0.1:49654"))
		service.Init()

		AuthClient = auth.NewAuthService("auth", service.Client())
		go func() {
			if err := service.Run(); err != nil {
				logger.Error("service auth run failed, err: %s", err.Error())
			}
		}()
	})
}
