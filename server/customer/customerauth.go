package customer

import (
    "context"
    "github.com/gin-gonic/gin"
    "time"
    "customermanager-go/api"
    "customermanager-go/server/common/db"
    "customermanager-go/server/common/db/dao"
    "customermanager-go/server/common/db/po"
    "customermanager-go/server/common/logger"
    "customermanager-go/server/common/resultcode"
    "customermanager-go/server/common/utils"
)

func Login(ctx context.Context, c *gin.Context) (api.BaseResponse, error) {
    ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*5))
    defer cancel()

    req := &api.LoginRequest{}
    if err := c.ShouldBindJSON(req); err != nil {
        return api.BaseResponse{
            Code:    resultcode.RequestIllegal,
            Message: resultcode.ResultMessage[resultcode.RequestIllegal],
            Data:    nil,
        }, nil
    }

    if req.UserName == "" || req.Password == "" {
        return api.BaseResponse{
            Code:    resultcode.UserNameOrPasswordError,
            Message: resultcode.ResultMessage[resultcode.UserNameOrPasswordError],
            Data:    nil,
        }, nil
    }

    session := db.Engine.NewSession()
    defer session.Close()

    customer, err := dao.Operator.QueryCustomerByName(session, req.UserName)
    if err != nil {
        logger.Error("query customer by name from db error, name: %s, error: %s", req.UserName, err.Error())
        return api.BaseResponse{
            Code:    resultcode.SystemInternalException,
            Message: resultcode.ResultMessage[resultcode.SystemInternalException],
            Data:    nil,
        }, nil
    }
    if customer == nil {
        return api.BaseResponse{
            Code:    resultcode.UserNameOrPasswordError,
            Message: resultcode.ResultMessage[resultcode.UserNameOrPasswordError],
            Data:    nil,
        }, nil
    }

    login, err := dao.Operator.QueryLoginByNameAndPwd(session, customer.Customerid, req.Password)
    if err != nil {
        logger.Error("query login by id from db error, id: %s, error: %s", customer.Customerid, err.Error())
        return api.BaseResponse{
            Code:    resultcode.SystemInternalException,
            Message: resultcode.ResultMessage[resultcode.SystemInternalException],
            Data:    nil,
        }, nil
    }
    if login == nil {
        return api.BaseResponse{
            Code:    resultcode.UserNameOrPasswordError,
            Message: resultcode.ResultMessage[resultcode.UserNameOrPasswordError],
            Data:    nil,
        }, nil
    }

    token := utils.NewUUID()
    newLogin := &po.LoginAuth{
        Customerid: customer.Customerid,
        Token:      token,
        Updatetime: time.Now(),
    }
    dao.Operator.AddLoginAuth(session, newLogin)

    logger.Info("user %s login success", req.UserName)

    dataJson := map[string]string{"customerId": customer.Customerid, "token": token}
    return api.BaseResponse{
        Code:    resultcode.Success,
        Message: resultcode.ResultMessage[resultcode.Success],
        Data:    dataJson,
    }, nil
}

func Logout(ctx context.Context, c *gin.Context) (api.BaseResponse, error) {
    return api.BaseResponse{Code: resultcode.Success, Message: resultcode.ResultMessage[resultcode.Success]}, nil
}

func Register(ctx context.Context, c *gin.Context) (api.BaseResponse, error) {
    return api.BaseResponse{Code: resultcode.Success, Message: resultcode.ResultMessage[resultcode.Success]}, nil
}
