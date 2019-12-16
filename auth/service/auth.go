package service

import (
	"customermanager-go/common/logger"
	"customermanager-go/common/proto/rpc/auth"
	"golang.org/x/net/context"
)

type AuthService struct{}

func (a *AuthService) Auth(ctx context.Context, req *auth.AuthRequest, response *auth.AuthResponse) error {
	logger.Info("receive request: %+v", req)
	response.Code = "0"
	response.Message = "success"
	return nil
}
