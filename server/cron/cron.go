package cron

import (
	error2 "customermanager-go/server/error"
	"customermanager-go/server/logger"
	"github.com/robfig/cron/v3"
	"sync"
)

type crontab struct {
	mutex sync.Mutex
	cron  *cron.Cron
	ids   map[string]cron.EntryID
}

func NewCrontab() *crontab {
	return &crontab{
		cron: cron.New(),
		ids:  make(map[string]cron.EntryID),
	}
}

func (c *crontab) Start() {
	c.cron.Start()
}

func (c *crontab) Stop() {
	c.cron.Stop()
}

func (c *crontab) Add(id string, cronExp string, job interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.ids[id]; ok {
		logger.Warn("cron is exist, id: %s", id)
		return nil
	}

	var entryId cron.EntryID
	var err error
	switch job.(type) {
	case cron.Job:
		entryId, err = c.cron.AddJob(cronExp, job.(cron.Job))
	case func():
		entryId, err = c.cron.AddFunc(cronExp, job.(func()))
	default:
		return error2.BaseError{
			Message: "job type not support",
		}
	}

	if err != nil {
		logger.Error("add job error, id: %s, error: %s", id, err.Error())
		return err
	}
	c.ids[id] = entryId
	logger.Info("add cron %s success", id)
	return nil
}

func (c *crontab) del(id string) {
	if _, ok := c.ids[id]; !ok {
		logger.Warn("cron %s is not exist", id)
		return
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entryId, ok := c.ids[id]
	if !ok {
		logger.Warn("cron %s is not exist", id)
		return
	}

	c.cron.Remove(entryId)
	delete(c.ids, id)
	logger.Info("delete cron %s success", id)
}
