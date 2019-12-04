package dao

import (
	"customermanager-go/server/common/constant"
	"customermanager-go/server/common/db/po"
	"customermanager-go/server/common/logger"
	"github.com/go-xorm/xorm"
)

var Operator daoImpl

type daoImpl struct {
}

type daoInter interface {
	// t_customer
	QueryCustomerByName(session *xorm.Session, customerName string) (*po.Customer, error)

	// t_login
	QueryLoginByNameAndPwd(session *xorm.Session, customerId string, password string) (*po.Login, error)

	// t_login_auth
	AddLoginAuth(session *xorm.Session, loginAuth *po.LoginAuth) error
}

func (d *daoImpl) QueryCustomerByName(session *xorm.Session, customerName string) (*po.Customer, error) {
	customer := &po.Customer{
		Customername: customerName,
	}

	session.Where("status = ?", constant.NORMAL)
	if _, err := session.Get(customer); err != nil {
		logger.Error("execute sql error, err: %s", err.Error())
		return nil, err
	}

	return customer, nil
}

func (d *daoImpl) QueryLoginByNameAndPwd(session *xorm.Session, customerId string, password string) (*po.Login, error) {
	login := &po.Login{
		Customerid: customerId,
		Password:   password,
	}
	if _, err := session.Get(login); err != nil {
		logger.Error("execute sql error, err: %s", err.Error())
		return nil, err
	}

	return login, nil
}

func (d *daoImpl) AddLoginAuth(session *xorm.Session, loginAuth *po.LoginAuth) error {
	if _, err := session.Insert(loginAuth); err != nil {
		logger.Error("execute sql error, err: %s", err.Error())
		return err
	}

	return nil
}
