package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"

	"../define"
	"../network"
	"../utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
)

// LogEngine 日志数据库
var LogEngine *xorm.Engine

// GameEngine 游戏数据库
var GameEngine *xorm.Engine

var logger *log.Logger

// Processor 处理器
type Processor struct {
	server *network.Server // 服务器
}

// OnMessage 收到消息
func (p *Processor) OnMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) error {
	ret := p.OnMessageEx(conn, mcmd, scmd, data)

	// 必须回复消息
	if ret == nil {
		return define.ErrSuccess
	}

	// 错误直接回复
	if err, ok := ret.(error); ok {
		return err
	}

	// 实现快捷回复消息
	return network.SendJSONMessage(conn, mcmd, scmd, ret)
}

// OnMessageEx 收到消息
func (p *Processor) OnMessageEx(conn net.Conn, mcmd uint16, scmd uint16, data []byte) interface{} {
	logger.WithFields(log.Fields{
		"mcmd": mcmd,
		"scmd": scmd,
		"data": string(data),
	}).Info("db OnMessage")

	switch mcmd {
	case define.DBCommon:
		return p.OnMainCommon(conn, scmd, data)
	}

	return define.ErrUnknownMainCmd
}

// OnMainCommon 通用主命令
func (p *Processor) OnMainCommon(conn net.Conn, scmd uint16, data []byte) interface{} {
	switch scmd {
	case define.DBFastRegister:
		return p.OnSubFastRegister(conn, data)
	case define.DBFastLogin:
		return p.OnSubFastLogin(conn, data)
	case define.DBChangeTreasure:
		return p.OnSubChangeTreasure(conn, data)
	}

	return define.ErrUnknownSubCmd
}

// ChangeUserTreasure 改变用户财富
func (p *Processor) ChangeUserTreasure(userId int, score int64, varScore int64, diamond int64, varDiamond int64, changeType int) error {
	var userTreasure UserTreasure
	// 当前分数钻石
	if score < 0 || diamond < 0 {
		if _, err := GameEngine.Where("user_id=?", userId).Get(&userTreasure); err != nil {
			return err
		}
	}

	// 更新分数钻石
	if _, err := GameEngine.Exec("UPDATE user_treasure SET user_score = user_score + ?, user_diamond = user_diamond + ? WHERE user_id = ?", varScore, varDiamond, userId); err != nil {
		return err
	}

	// 记录财富日志
	if _, err := LogEngine.Exec(fmt.Sprintf("INSERT INTO user_treasure_log_%d (user_id, cur_score, var_score, cur_diamond, var_diamond, change_type) VALUES (?, ?, ?, ?, ?, ?)", utils.Date()), userId, userTreasure.UserScore, varScore, userTreasure.UserDiamond, varDiamond, changeType); err != nil {
		return err
	}

	return nil
}

// OnSubFastRegister 快速注册子命令
func (p *Processor) OnSubFastRegister(conn net.Conn, data []byte) interface{} {
	fastRegister := &define.FastRegister{}
	replyFastRegister := &define.ReplyFastRegister{}

	if err := json.Unmarshal(data, fastRegister); err != nil {
		return err
	}

	// 查询用户信息
	userInfo := new(UserInformation)
	_, err := GameEngine.Where("user_account = ?", fastRegister.Account).Get(&userInfo)

	if err == sql.ErrNoRows {
		// 新建用户信息
		newUser := User{
			Account:  fastRegister.Account,
			Password: fastRegister.Password,
			Name:     fastRegister.Name,
			Icon:     fastRegister.Icon,
			Ip:       fastRegister.IP,
			Gender:   uint8(fastRegister.Gender),
		}
		_, err := GameEngine.Insert(&newUser)

		if err != nil {
			return err
		}

		// 新建用户附加信息
		userInfo.UserId = newUser.Id
		userInfo.UserDiamond = 0
		userInfo.UserScore = 0
		userInfo.UserLevel = 1
		userInfo.UserName = newUser.Name

		_,err =GameEngine.Insert(userInfo)
		if err != nil {
			return err
		}

		replyFastRegister.UserID = int(newUser.Id)

		// 插入用户财富
		newUserTreasure := UserTreasure{
			UserId:      newUser.Id,
			UserScore:   0,
			UserDiamond: 0,
		}
		_,err = GameEngine.Insert(newUserTreasure)
		if err != nil {
			return err
		}

		// 用户初始分数钻石
		var score, diamond int64

		// @todo 配置信息建议从redis中读取
		//if err := GameEngine.QueryRow(`SELECT Content FROM game_config WHERE Title = "InitScore"`).Scan(&score); err != nil {
		//	return err
		//}
		//
		//if err := GameEngine.QueryRow(`SELECT Content FROM game_config WHERE Title = "InitDiamond"`).Scan(&diamond); err != nil {
		//	return err
		//}

		score = 10
		diamond = 10
		replyFastRegister.UserScore = score
		replyFastRegister.UserDiamond = diamond

		// 用户财富变化
		if err := p.ChangeUserTreasure(int(newUser.Id), 0, score, 0, diamond, define.ChangeTypeRegister); err != nil {
			return err
		}

		// 初始用户等级
		replyFastRegister.UserLevel = userInfo.UserLevel
	} else if err != nil {
		return err
	}

	// 总是更新这些字段
	userInfo.UserName = fastRegister.Name
	userInfo.UserIcon = fastRegister.Icon
	userInfo.UserGender = fastRegister.Gender
	_,err = GameEngine.Where("user_id = ?", userInfo.UserId).Update(userInfo)
	if err != nil {
		return err
	}

	return replyFastRegister
}

// OnSubFastLogin 快速登陆子命令
func (p *Processor) OnSubFastLogin(conn net.Conn, data []byte) interface{} {
	var id int
	//replyFastLogin := &define.ReplyFastLogin{}

	if err := json.Unmarshal(data, &id); err != nil {
		return err
	}

	// 查询用户信息
	userInfo := new(UserInformation);
	_, err := GameEngine.Where("user_id=?", id).Get(&userInfo)

	if err != nil {
		return err
	}

	return userInfo
}

// OnSubChangeTreasure 改变财富
func (p *Processor) OnSubChangeTreasure(conn net.Conn, data []byte) interface{} {
	notifyTreasure := &define.NotifyTreasure{}

	if err := json.Unmarshal(data, notifyTreasure); err != nil {
		return err
	}

	// 用户财富变化
	return p.ChangeUserTreasure(notifyTreasure.UserID,
		-1, notifyTreasure.VarScore,
		-1, notifyTreasure.VarDiamond,
		notifyTreasure.ChangeType)
}

// OnClose 连接关闭
func (p *Processor) OnClose(conn net.Conn) {

}

// OnClientMessage 客户端收到消息
func (p *Processor) OnClientMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) {
	// nothing to do
}

// OnClientConnect 客户端连接成功
func (p *Processor) OnClientConnect(conn net.Conn) {
	// nothing to do
}

// NewProcessor 创建处理器
func NewProcessor(server *network.Server, config *define.GConfig) *Processor {
	logger = utils.GetLogger("db")
	var err error

	// todo SetMaxOpenConns, SetMaxIdleConns

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
		return nil
	}

	if err = GameEngine.Ping(); err != nil {
		logger.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Ping game")
		return nil
	}

	//初始化游戏数据库
	GameDbInit(config)
	logger.Info("db init finished")

	return &Processor{
		server: server,
	}
}

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
