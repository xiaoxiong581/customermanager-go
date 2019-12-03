package http

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    uuid "github.com/satori/go.uuid"
    "io/ioutil"
    "net/http"
    "time"
    "customermanager-go/api"
    "customermanager-go/server/common/logger"
    "customermanager-go/server/common/resultcode"
    "customermanager-go/server/common/utils"
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

    sessionID := uuid.NewV4().String()
    startTime := time.Now()
    reqBody := make(map[string]interface{})
    reqBodyBytes, _ := ioutil.ReadAll(c.Request.Body)
    c.Request.Body = ioutil.NopCloser(bytes.NewReader(reqBodyBytes))
    url := c.Request.URL.String()
    err := json.NewDecoder(bytes.NewReader(reqBodyBytes)).Decode(&reqBody)

    sessionIDData, ok := reqBody[SessionKey]
    if ok {
        _sessionID, _ok := sessionIDData.(string)
        if _ok && _sessionID != "" {
            sessionID = _sessionID
        }
    }
    ctx := utils.NewContext(sessionID, "")
    logger.Info("[server] get request, time: %s, method: %s, path: %s, url: %s, reqBody: %s, decodeRequestBodyResult: %s", startTime, method, path, url, reqBody, err)

    c.Keys = map[string]interface{}{}
    c.Keys["reqBodyBytes"] = reqBodyBytes
    c.Keys["reqBody"] = reqBody

    res, error := handlerFunc(ctx, c)
    handleEndTime := time.Now()
    logger.Info("[server] handle end, method: %s, path: %s, url: %s, start_time: %s, end_time: %s, func_handle_time: %d, handle_time_unit: %s, response: %s, error: %s, http_status_code: %d", method, path, url, startTime.Format(TimeFormat), handleEndTime.Format(TimeFormat), (handleEndTime.UnixNano()-startTime.UnixNano())/1000000, "ms", res, error, http.StatusOK)

    if error != nil {
        logger.Info("handle error, error: %+v", error)

        c.AbortWithStatusJSON(http.StatusOK, &api.BaseResponse{
            Code:    resultcode.SystemInternalException,
            Message: resultcode.ResultMessage[resultcode.SystemInternalException],
        })

        return
    }

    logger.Info("handle success")
    c.JSON(http.StatusOK, res)
}
