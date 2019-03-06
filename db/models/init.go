package models

import (
	"../../define"
	"../../utils"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
)

var GameEngine *xorm.Engine

var LogEngine *xorm.Engine

var logger *log.Logger

// 初始化数据库
func InitDB(config *define.GConfig) error  {
	var err error
	// 初始化日志组件
	if logger == nil {
		logger = utils.GetLogger("db")
	}

	if LogEngine, err = xorm.NewEngine("mysql", "root:123456@/game_log"); err != nil {
		logger.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Open log db")
		return nil
	}

	if err = LogEngine.Ping(); err != nil {
		logger.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Ping log")
		return nil
	}

	if GameEngine, err = xorm.NewEngine("mysql", "root:123456@/game_db"); err != nil {
		logger.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Open game")
		return err
	}

	if err = GameEngine.Ping(); err != nil {
		logger.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Ping game")
		return err
	}

	GameDbInit(config)
	return  nil
}

func CloseEngine() {
	GameEngine.Close()
	LogEngine.Close()
}
// 获得Logger
func GetDBLogger() *log.Logger  {
	return logger
}

func DBLogError(msg string,err error) {
	logger.WithFields(log.Fields{
		"error": err.Error(),
	}).Error(msg)
}

// 游戏数据库初始化
func GameDbInit(config *define.GConfig) {
	// 设置表前缀
	tableMapper := core.NewPrefixMapper(core.GonicMapper{}, config.DB.GameTablePrefix)
	GameEngine.SetTableMapper(tableMapper)

	//初始化用户表
	//if tbExist, err := GameEngine.IsTableExist(&User{}); err != nil && !tbExist {
	//
	//}
	err := GameEngine.Sync2(new(User), new(UserTreasure), new(UserInformation))
	if err != nil {
		logger.WithFields(log.Fields{
			"db_error" : err.Error(),
		}).Error("sync table failed")
	}
}
