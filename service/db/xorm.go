package db

import (
	_ "github.com/go-sql-driver/mysql"
	"learn_together/initer"
	"os"
	"time"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

var Engine *xorm.Engine

func InitXorm(config *initer.Config) {
	var err error
	Engine, err = xorm.NewEngine("mysql", config.Mysql.Connection)
	if err != nil {
		panic("mysql init error" + err.Error())
	}
	Engine.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
	//日志
	logFile, err := os.OpenFile(config.LogPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic("log file can't initer")
	}
	logger := log.NewSimpleLogger(logFile)
	logger.SetLevel(log.LOG_INFO)
	Engine.SetLogger(log.NewLoggerAdapter(logger))
	Engine.ShowSQL(true)
}
