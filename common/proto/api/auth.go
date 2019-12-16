package api

import (
	"customermanager-go/common/proto/rpc/auth"
	"golang.org/x/net/context"
)

type AuthApi interface {
	Auth(ctx context.Context, req *auth.AuthRequest, response *auth.AuthResponse) error
}
