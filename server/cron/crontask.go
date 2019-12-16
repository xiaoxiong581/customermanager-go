package cron

import (
	cron2 "customermanager-go/common/cron"
	db2 "customermanager-go/common/db"
	"customermanager-go/common/logger"
	"customermanager-go/server/db/dao"
	"flag"
	"time"
)

const TIME_FORMAT = "2006-01-02 15:04:05.999"

var existTime = flag.Duration("tokenExistTime", 10, "token exist time")

func StartCron() {
	cron := cron2.NewCrontab()
	if err := cron.Add("cleanExpireLoginAuthCron", "@every 1m", cleanExpireLoginAuthCron); err != nil {
		logger.Error("add cleanExpireLoginAuthCron fail, error: %s", err.Error())
	}

	cron.Start()
}

func cleanExpireLoginAuthCron() {
	session := db2.Engine.NewSession()
	defer session.Close()

	expireTime := time.Now().Add(-time.Minute * *existTime)
	logger.Info("begin to delete expire login auth, expireTime: %s", expireTime.Format(TIME_FORMAT))
	dao.Operator.DeleteExpireLoginAuth(session, expireTime)
}
