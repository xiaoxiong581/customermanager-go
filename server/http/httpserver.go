package http

import (
	"bytes"
	db2 "customermanager-go/common/db"
	error2 "customermanager-go/common/error"
	"customermanager-go/common/logger"
	"customermanager-go/server/db/dao"
	"customermanager-go/server/db/po"
	"customermanager-go/server/resultcode"
	"customermanager-go/server/service/api"
	"customermanager-go/server/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	SessionKey = "request_uuid"
	TimeFormat = "2006-01-02T15:04:05.999999999-07:00"
)

type ServerConfig struct {
	IP   string
	Port int
}

func NewHttpServer(config *ServerConfig) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.IP, config.Port),
		Handler: newRouter(),
	}
}

func newRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	for index := range Routes {
		route := Routes[index]
		path := route.Pattern
		method := route.Method

		if method == POST {
			router.POST(path, func(c *gin.Context) {
				loader(c, route)
			})
		}
		if method == GET {
			router.GET(path, func(c *gin.Context) {
				loader(c, route)
			})
		}
		if method == PUT {
			router.PUT(path, func(c *gin.Context) {
				loader(c, route)
			})
		}
		if method == DELETE {
			router.DELETE(path, func(c *gin.Context) {
				loader(c, route)
			})
		}
		if method == PATCH {
			router.PATCH(path, func(c *gin.Context) {
				loader(c, route)
			})
		}
	}

	return router
}

func loader(c *gin.Context, route Router) {
	handlerFunc := route.Func
	path := route.Pattern
	method := route.Method
	url := c.Request.URL.String()

	if _, has := AuthExcludeUrls[url]; !has {
		authSuc, err := userAuth(c.GetHeader("auth_customerId"), c.GetHeader("auth_token"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, &api.BaseResponse{
				Code:    resultcode.SystemInternalException,
				Message: resultcode.ResultMessage[resultcode.SystemInternalException],
			})
			return
		}
		if !authSuc {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &api.BaseResponse{
				Code:    resultcode.UserAuthFail,
				Message: resultcode.ResultMessage[resultcode.UserAuthFail],
			})
			return
		}
	}

	sessionID := uuid.NewV4().String()
	startTime := time.Now()
	reqBody := make(map[string]interface{})
	reqBodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	c.Request.Body = ioutil.NopCloser(bytes.NewReader(reqBodyBytes))
	if err := json.NewDecoder(bytes.NewReader(reqBodyBytes)).Decode(&reqBody); err != nil {
		logger.Error("decode request body error, err: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusOK, &api.BaseResponse{
			Code:    resultcode.SystemInternalException,
			Message: resultcode.ResultMessage[resultcode.SystemInternalException],
		})
		return
	}

	sessionIDData, ok := reqBody[SessionKey]
	if ok {
		_sessionID, _ok := sessionIDData.(string)
		if _ok && _sessionID != "" {
			sessionID = _sessionID
		}
	}
	ctx := utils.NewContext(sessionID, "")
	logger.Info("[server] get request, time: %s, method: %s, path: %s, url: %s, reqBody: %s", startTime, method, path, url, reqBody)

	res, err := handlerFunc(ctx, c)
	handleEndTime := time.Now()
	logger.Info("[server] handle end, method: %s, path: %s, url: %s, start_time: %s, end_time: %s, handle_time: %dms, response: %+v", method, path, url, startTime.Format(TimeFormat), handleEndTime.Format(TimeFormat), (handleEndTime.UnixNano()-startTime.UnixNano())/1000000, res)

	if err != nil {
		logger.Error("handle error, error: %s", err.Error())
		code := resultcode.SystemInternalException
		message := resultcode.ResultMessage[resultcode.SystemInternalException]
		if _, ok := err.(error2.BaseError); ok {
			code = err.(error2.BaseError).Code
			message = err.(error2.BaseError).Message
		}

		c.AbortWithStatusJSON(http.StatusOK, &api.BaseResponse{
			Code:    code,
			Message: message,
		})
		return
	}

	logger.Info("handle success")
	c.JSON(http.StatusOK, res)
}

func userAuth(customerId string, token string) (bool, error) {
	if customerId == "" || token == "" {
		return false, nil
	}

	session := db2.Engine.NewSession()
	defer session.Close()

	loginAuth := &po.LoginAuth{
		Customerid: customerId,
		Token:      token,
		Updatetime: time.Now(),
	}

	updateNum, err := dao.Operator.UpdateTimeByAuth(session, loginAuth)
	if err != nil {
		logger.Error("query login auth failed, error: %s", err.Error())
		return false, err
	}

	return updateNum > 0, nil
}
