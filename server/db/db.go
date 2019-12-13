package db

import (
	"customermanager-go/server/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"strconv"
	"xorm.io/core"
)

var Engine *xorm.Engine

const (
	DEFAULT_ALLOWED_PACKET = 16 * 1024 * 1024
)

func InitEngine(dbString string, maxAllowedPacket int) error {
	var err error
	if maxAllowedPacket < DEFAULT_ALLOWED_PACKET {
		maxAllowedPacket = DEFAULT_ALLOWED_PACKET
	}

	dbSourceName := dbString + "&maxAllowedPacket=" + strconv.Itoa(maxAllowedPacket)
	logger.Info("init db, url: %s", dbString)
	Engine, err = xorm.NewEngine("mysql", dbSourceName)
	if err != nil {
		logger.Error("init db error, error: %s", err.Error())
		return err
	}

	Engine.SetTableMapper(core.NewPrefixMapper(core.SnakeMapper{}, "t_"))
	Engine.ShowSQL(true)

	go func() {
		if err := Engine.Ping(); err != nil {
			logger.Error("ping db error, error: %s", err.Error())
		}
	}()

	logger.Info("init db success")
	return nil
}
