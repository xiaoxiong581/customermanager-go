package po

import (
	"time"
)

/**
 * generate by xorm tool, like: xorm reverse -s mysql root:xiaoxiong581@tcp\(192.168.137.106:3306\)/merchant?charset=utf8 templates/goxorm
 */
type Customer struct {
	Customerid   string    `xorm:"not null pk comment('用户Id') VARCHAR(64)"`
	Customername string    `xorm:"not null comment('用户名') unique VARCHAR(128)"`
	Status       int       `xorm:"not null default 0 comment('用户状态 0: 正常; 1: 锁定; 99: 删除') INT(11)"`
	Createtime   time.Time `xorm:"not null comment('创建时间') DATETIME"`
	Updatetime   time.Time `xorm:"not null comment('更新时间') DATETIME"`
	Email        string    `xorm:"not null comment('邮箱') VARCHAR(128)"`
}

type Login struct {
	Customerid string    `xorm:"not null pk comment('用户Id') VARCHAR(64)"`
	Password   string    `xorm:"not null comment('用户密码') VARCHAR(64)"`
	Createtime time.Time `xorm:"not null comment('创建时间') DATETIME"`
	Updatetime time.Time `xorm:"not null comment('更新时间') DATETIME"`
}

type LoginAuth struct {
	Customerid string    `xorm:"not null pk comment('用户Id') unique(unique_customerId_token) VARCHAR(64)"`
	Token      string    `xorm:"not null pk comment('鉴权token') unique(unique_customerId_token) VARCHAR(64)"`
	Updatetime time.Time `xorm:"not null comment('更新时间') DATETIME"`
}
