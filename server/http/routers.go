package http

import (
	"context"
	"customermanager-go/server/service"
	"customermanager-go/server/service/api"
	"github.com/gin-gonic/gin"
)

const (
	POST   = "post"
	GET    = "get"
	PUT    = "put"
	DELETE = "delete"
	PATCH  = "patch"

	BASE_URL     = "/customermanager/v1/customer"
	LOGIN_URL    = BASE_URL + "/auth/login"
	LOGOUT_URL   = BASE_URL + "/auth/logout"
	REGISTER_URL = BASE_URL + "/auth/register"
)

var AuthExcludeUrls = map[string]bool{
	LOGIN_URL:    true,
	LOGOUT_URL:   true,
	REGISTER_URL: true,
}

type httpHandler func(ctx context.Context, c *gin.Context) (api.BaseResponse, error)

type Router struct {
	Method  string
	Pattern string
	Func    httpHandler
}

var (
	Routes = []Router{
		{Method: POST, Pattern: LOGIN_URL, Func: service.Login},
		{Method: POST, Pattern: LOGOUT_URL, Func: service.Logout},
		{Method: POST, Pattern: REGISTER_URL, Func: service.Register},
	}
)
