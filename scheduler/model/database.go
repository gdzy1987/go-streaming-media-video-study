package model

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go-streaming-media-video-study/config"
	"go-streaming-media-video-study/logger"
)

var (
	dbConn *sql.DB
	err    error
)

func InitMysql() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:3306)/video_server?charset=utf8", config.DefaultConfig.MysqlUser,
		config.DefaultConfig.MysqlPassword, config.DefaultConfig.MysqlIP)

	if dbConn, err = sql.Open("mysql", dataSourceName); err != nil {
		logger.Info("connect mysql error:\t", dataSourceName, err)
		panic(err.Error())
	}
}
