package db

import (
	"encoding/json"
	"fmt"
	"net"

	"../define"
	"../network"
	"../utils"
	"./models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
)

// LogEngine 日志数据库
var LogEngine *xorm.Engine

// GameEngine 游戏数据库
var GameEngine *xorm.Engine

//  日志
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
	case define.DBRegisterCheck:
		return p.OnSubRegisterCheck(conn, data)
	}

	return define.ErrUnknownSubCmd
}

// ChangeUserTreasure 改变用户财富
func (p *Processor) ChangeUserTreasure(userId int, varScore int64, varDiamond int64, changeType int) error {
	var userTreasure *models.UserTreasure
	userTreasureModel := models.NewUserTreasure()
	// 当前分数钻石
	userTreasure,err := userTreasureModel.GetTreasureByUserId(uint32(userId))
	if err != nil {
		return err
	}

	// 更新分数钻石
	if err := userTreasureModel.UpdateUserTreasure(uint32(userId),varScore, varDiamond); err !=nil {
		return err
	}

	// 记录财富日志
	if _, err := models.LogEngine.Exec(fmt.Sprintf("INSERT INTO user_treasure_log_%d (user_id, cur_score, var_score, cur_diamond, var_diamond, change_type) VALUES (?, ?, ?, ?, ?, ?)", utils.Date()), userId, userTreasure.UserScore, varScore, userTreasure.UserDiamond, varDiamond, changeType); err != nil {
		return err
	}

	return nil
}

func (p *Processor) OnSubRegisterCheck(conn net.Conn, data []byte)error {
	registerCheck := &define.FastRegisterCheck{}
	if err := json.Unmarshal(data, registerCheck); err != nil {
		logger.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("OnSubRegisterCheck json unmarshal")
		return err
	}
	userInfoModel := models.NewUserInfo()
	userInfo,err := userInfoModel.GetInfoByAccount(registerCheck.Account)
	if err != nil {
		return err
	}else if userInfo != nil {
		return define.ErrRegisterAccountExist
	}
	userInfo,err = userInfoModel.GetInfoByName(registerCheck.Name)
	if err != nil {
		return err
	}else if userInfo != nil {
		return define.ErrRegisterNameExist
	}
	return nil
}

// OnSubFastRegister 快速注册子命令
func (p *Processor) OnSubFastRegister(conn net.Conn, data []byte) interface{} {
	fastRegister := &define.FastRegister{}
	replyFastRegister := &define.ReplyFastRegister{}

	if err := json.Unmarshal(data, fastRegister); err != nil {
		logger.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("OnSubFastRegister json unmarshal")
		return err
	}

	logger.WithFields(log.Fields{
		"RegisterInfo" : fastRegister,
	}).Info("FastRegister data")

	// 查询用户信息
	userInfoModel := new(models.UserInformation)
	userInfoModel,err := userInfoModel.GetInfoByName(fastRegister.Name)
	if userInfoModel != nil {
		return define.ErrRegisterNameExist
	}

	userInfo,err := userInfoModel.GetInfoByAccount(fastRegister.Account)

	logger.WithFields(log.Fields{
		"UserInfo" : userInfo,
	}).Info("query user info data")

	if userInfo == nil {
		// 新建用户信息
		userModel := models.NewUser()
		newUser,err := userModel.RegisterUser(fastRegister)
		if err != nil {
			return err
		}

		logger.WithFields(log.Fields{
			"user" : newUser,
		}).Info("add new user")

		// 新建用户附加信息
		userInfoModel := models.NewUserInfo()
		userInfo,err := userInfoModel.AddUserInfoByUser(newUser)
		replyFastRegister.UserID = uint32(newUser.Id)

		// 插入用户财富
		userTreasureModel := models.NewUserTreasure()
		userTreasureModel.InitUserTreasure(newUser.Id)

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
		if err := p.ChangeUserTreasure(int(newUser.Id), score, diamond, define.ChangeTypeRegister); err != nil {
			return err
		}

		// 初始用户等级
		replyFastRegister.UserLevel = userInfo.UserLevel
	} else if err != nil {
		logger.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("db query user info by account error")
		return err
	} else {//新增逻辑账号存在返回错误信息
		return define.ErrRegisterAccountExist
	}

	return replyFastRegister
}

// OnSubFastLogin 快速登陆子命令
func (p *Processor) OnSubFastLogin(conn net.Conn, data []byte) interface{} {
	var id uint32

	if err := json.Unmarshal(data, &id); err != nil {
		return err
	}

	// 查询用户信息
	userInfoModel := models.NewUserInfo()
	userInfo,err :=userInfoModel.GetInfoByUserId(id)
	if err != nil {
		return err
	}
	replyFastLogin := &define.ReplyFastLogin{}
	replyFastLogin.UserID = userInfo.UserId
	replyFastLogin.UserName=userInfo.UserName
	replyFastLogin.IsRobot= (userInfo.IsRobot > 0)
	replyFastLogin.UserGender=userInfo.UserGender
	replyFastLogin.UserIcon=userInfo.UserIcon
	replyFastLogin.UserScore=userInfo.UserScore
	replyFastLogin.UserDiamond=userInfo.UserDiamond
	replyFastLogin.BindPhone = userInfo.BindPhone

	return replyFastLogin
}

// OnSubChangeTreasure 改变财富
func (p *Processor) OnSubChangeTreasure(conn net.Conn, data []byte) interface{} {
	notifyTreasure := &define.NotifyTreasure{}

	if err := json.Unmarshal(data, notifyTreasure); err != nil {
		return err
	}

	// 用户财富变化
	return p.ChangeUserTreasure(notifyTreasure.UserID,
		 notifyTreasure.VarScore,
		 notifyTreasure.VarDiamond,
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

	// todo SetMaxOpenConns, SetMaxIdleConns
	models.InitDB(config)
	logger = models.GetDBLogger()
	//初始化游戏数据库
	logger.Info("db init finished")

	return &Processor{
		server: server,
	}
}

