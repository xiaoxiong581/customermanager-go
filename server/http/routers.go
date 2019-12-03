package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"customermanager-go/api"
	"customermanager-go/server/customer"
)

const (
	POST   = "post"
	GET    = "get"
	PUT    = "put"
	DELETE = "delete"
	PATCH  = "patch"

	BASE_URL = "/customermanager/v1/customer"
)

type httpHandler func(ctx context.Context, c *gin.Context) (api.BaseResponse, error)

type Router struct {
	Method  string
	Pattern string
	Func    httpHandler
}

var (
	Routes = []Router{
		{Method: POST, Pattern: BASE_URL + "/auth/login", Func: customer.Login},
		{Method: POST, Pattern: BASE_URL + "/auth/logout", Func: customer.Logout},
		{Method: POST, Pattern: BASE_URL + "/auth/register", Func: customer.Register},
	}
)
