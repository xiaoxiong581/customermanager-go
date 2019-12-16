package service

import (
	"context"
	"customermanager-go/common/db"
	"customermanager-go/common/logger"
	"customermanager-go/server/constant"
	"customermanager-go/server/db/dao"
	"customermanager-go/server/db/po"
	"customermanager-go/server/resultcode"
	"customermanager-go/server/service/api"
	"customermanager-go/server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func Login(ctx context.Context, c *gin.Context) (api.BaseResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*5))
	defer cancel()

	req := &api.LoginRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		return api.BaseResponse{
			Code:    resultcode.RequestIllegal,
			Message: fmt.Sprintf(resultcode.ResultMessage[resultcode.RequestIllegal], err.Error()),
			Data:    nil,
		}, nil
	}

	/* // test go-micro framework
	   authRequest := &auth.AuthRequest{
	       CustomerId: req.UserName,
	       Token:      req.Password,
	   }
	   authResponse, err := rpc.AuthClient.Auth(ctx, authRequest)
	   if err != nil {
	       logger.Error("send rpc to auth failed, %s", err.Error())
	       return api.BaseResponse{}, err
	   }
	   logger.Info("receive auth reponse, %+v", authResponse)*/

	session := db.Engine.NewSession()
	defer session.Close()

	customer, has, err := dao.Operator.QueryCustomerByName(session, req.UserName)
	if err != nil {
		logger.Error("query customer by name from db error, name: %s, error: %s", req.UserName, err.Error())
		return api.BaseResponse{}, err
	}
	if !has {
		return api.BaseResponse{
			Code:    resultcode.UserNameOrPasswordError,
			Message: resultcode.ResultMessage[resultcode.UserNameOrPasswordError],
			Data:    nil,
		}, nil
	}
	if customer.Status != constant.NORMAL {
		return api.BaseResponse{
			Code:    resultcode.UserStatusError,
			Message: resultcode.ResultMessage[resultcode.UserStatusError],
			Data:    nil,
		}, nil
	}

	_, has, err = dao.Operator.QueryLoginByNameAndPwd(session, customer.Customerid, req.Password)
	if err != nil {
		logger.Error("query login by id from db error, id: %s, error: %s", customer.Customerid, err.Error())
		return api.BaseResponse{}, err
	}
	if !has {
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
	if err = dao.Operator.AddLoginAuth(session, newLogin); err != nil {
		logger.Error("add login auth to db error, customerName: %s", req.UserName)
		return api.BaseResponse{}, err
	}

	logger.Info("user %s login success", req.UserName)

	dataJson := map[string]string{"customerId": customer.Customerid, "token": token}
	return api.BaseResponse{
		Code:    resultcode.Success,
		Message: resultcode.ResultMessage[resultcode.Success],
		Data:    dataJson,
	}, nil
}

func Logout(ctx context.Context, c *gin.Context) (api.BaseResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*5))
	defer cancel()

	req := &api.LogoutRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		return api.BaseResponse{
			Code:    resultcode.RequestIllegal,
			Message: fmt.Sprintf(resultcode.ResultMessage[resultcode.RequestIllegal], err.Error()),
			Data:    nil,
		}, nil
	}

	session := db.Engine.NewSession()
	defer session.Close()

	loginAuth := &po.LoginAuth{
		Customerid: req.CustomerId,
		Token:      req.Token,
	}
	dao.Operator.DeleteLoginAuth(session, loginAuth)

	return api.BaseResponse{Code: resultcode.Success, Message: resultcode.ResultMessage[resultcode.Success]}, nil
}

func Register(ctx context.Context, c *gin.Context) (api.BaseResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*5))
	defer cancel()

	req := &api.RegisterRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		return api.BaseResponse{
			Code:    resultcode.RequestIllegal,
			Message: fmt.Sprintf(resultcode.ResultMessage[resultcode.RequestIllegal], err.Error()),
			Data:    nil,
		}, nil
	}

	session := db.Engine.NewSession()
	defer session.Close()

	_, has, err := dao.Operator.QueryCustomerByName(session, req.CustomerName)
	if err != nil {
		return api.BaseResponse{}, err
	}
	if has {
		return api.BaseResponse{
			Code:    resultcode.CustomerNameAlreadyExist,
			Message: fmt.Sprintf(resultcode.ResultMessage[resultcode.CustomerNameAlreadyExist], req.CustomerName),
			Data:    nil,
		}, nil
	}

	_, has, err = dao.Operator.QueryCustomerByEmail(session, req.Email)
	if err != nil {
		return api.BaseResponse{}, err
	}
	if has {
		return api.BaseResponse{
			Code:    resultcode.EmailAlreadyExist,
			Message: fmt.Sprintf(resultcode.ResultMessage[resultcode.EmailAlreadyExist], req.Email),
			Data:    nil,
		}, nil
	}

	customerId := utils.NewUUID()
	if err := session.Begin(); err != nil {
		return api.BaseResponse{}, err
	}
	customer := &po.Customer{
		Customerid:   customerId,
		Customername: req.CustomerName,
		Status:       constant.NORMAL,
		Createtime:   time.Now(),
		Updatetime:   time.Now(),
		Email:        req.Email,
	}
	err = dao.Operator.AddCustomer(session, customer)
	if err != nil {
		logger.Error("add customer info to db error, customerName: %s", req.CustomerName)
		return api.BaseResponse{}, err
	}

	login := &po.Login{
		Customerid: customerId,
		Password:   req.Password,
		Createtime: time.Now(),
		Updatetime: time.Now(),
	}
	err = dao.Operator.AddLogin(session, login)
	if err != nil {
		logger.Error("add login info to db error, customerName: %s", req.CustomerName)
		return api.BaseResponse{}, err
	}
	session.Commit()

	logger.Info("customerName %s register success", req.CustomerName)
	return api.BaseResponse{Code: resultcode.Success, Message: resultcode.ResultMessage[resultcode.Success]}, nil
}
