package cron

import (
	"customermanager-go/server/db"
	"customermanager-go/server/db/dao"
	"customermanager-go/server/logger"
	"flag"
	"time"
)

const TIME_FORMAT = "2006-01-02 15:04:05.999"

func StartCrontab() {
	crontab := NewCrontab()
	if err := crontab.Add("cleanExpireLoginAuthCron", "@every 1m", cleanExpireLoginAuthCron); err != nil {
		logger.Error("add cleanExpireLoginAuthCron fail, error: %s", err.Error())
	}

	crontab.Start()
}

func cleanExpireLoginAuthCron() {
	session := db.Engine.NewSession()
	defer session.Close()

	existTime := flag.Duration("tokenExistTime", 10, "token exist time")
	expireTime := time.Now().Add(-time.Minute * *existTime)
	logger.Info("begin to delete expire login auth, expireTime: %s", expireTime.Format(TIME_FORMAT))
	dao.Operator.DeleteExpireLoginAuth(session, expireTime)
}
