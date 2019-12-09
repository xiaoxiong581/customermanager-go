package dao

import (
	"customermanager-go/server/db/po"
	"customermanager-go/server/logger"
	"github.com/go-xorm/xorm"
	"time"
)

var Operator daoImpl

type daoImpl struct {
}

type daoInter interface {
	// t_customer
	AddCustomer(session *xorm.Session, customer *po.Customer) error
	QueryCustomerByName(session *xorm.Session, customerName string) (*po.Customer, bool, error)
	QueryCustomerByEmail(session *xorm.Session, email string) (*po.Customer, bool, error)

	// t_login
	AddLogin(session *xorm.Session, customer *po.Login) error
	QueryLoginByNameAndPwd(session *xorm.Session, customerId string, password string) (*po.Login, bool, error)

	// t_login_auth
	AddLoginAuth(session *xorm.Session, loginAuth *po.LoginAuth) error
	UpdateTimeByAuth(session *xorm.Session, loginAuth *po.LoginAuth) (bool, error)
	DeleteLoginAuth(session *xorm.Session, loginAuth *po.LoginAuth) error
	DeleteExpireLoginAuth(session *xorm.Session, expireTime time.Time) error
}

func (d *daoImpl) AddCustomer(session *xorm.Session, customer *po.Customer) error {
	if _, err := session.Insert(customer); err != nil {
		logger.Error("execute sql error, err: %s", err.Error())
		return err
	}

	return nil
}

func (d *daoImpl) QueryCustomerByName(session *xorm.Session, customerName string) (*po.Customer, bool, error) {
	customer := &po.Customer{
		Customername: customerName,
	}

	has, err := session.Get(customer)
	if err != nil {
		logger.Error("execute sql error, err: %s", err.Error())
		return nil, has, err
	}

	return customer, has, nil
}

func (d *daoImpl) QueryCustomerByEmail(session *xorm.Session, email string) (*po.Customer, bool, error) {
	customer := &po.Customer{
		Email: email,
	}

	has, err := session.Get(customer)
	if err != nil {
		logger.Error("execute sql error, err: %s", err.Error())
		return nil, has, err
	}

	return customer, has, nil
}

func (d *daoImpl) AddLogin(session *xorm.Session, login *po.Login) error {
	if _, err := session.Insert(login); err != nil {
		logger.Error("execute sql error, err: %s", err.Error())
		return err
	}

	return nil
}

func (d *daoImpl) QueryLoginByNameAndPwd(session *xorm.Session, customerId string, password string) (*po.Login, bool, error) {
	login := &po.Login{
		Customerid: customerId,
		Password:   password,
	}

	has, err := session.Get(login)
	if err != nil {
		logger.Error("execute sql error, err: %s", err.Error())
		return nil, has, err
	}

	return login, has, nil
}

func (d *daoImpl) AddLoginAuth(session *xorm.Session, loginAuth *po.LoginAuth) error {
	if _, err := session.Insert(loginAuth); err != nil {
		logger.Error("execute sql error, err: %s", err.Error())
		return err
	}

	return nil
}

func (d *daoImpl) UpdateTimeByAuth(session *xorm.Session, loginAuth *po.LoginAuth) (int64, error) {
	updateNum, err := session.Where("customerId = ? and token = ?", loginAuth.Customerid, loginAuth.Token).Cols("updateTime").Update(loginAuth)
	if err != nil {
		logger.Error("execute sql error, err: %s", err.Error())
		return 0, err
	}

	return updateNum, nil
}

func (d *daoImpl) DeleteLoginAuth(session *xorm.Session, loginAuth *po.LoginAuth) error {
	_, err := session.Delete(loginAuth)
	return err
}

func (d *daoImpl) DeleteExpireLoginAuth(session *xorm.Session, expireTime time.Time) error {
	_, err := session.Where("updateTime <= ?", expireTime).Delete(po.LoginAuth{})
	return err
}
