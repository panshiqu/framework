package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/network"
	"github.com/panshiqu/framework/utils"
)

// LOG 日志数据库
var LOG *sql.DB

// GAME 游戏数据库
var GAME *sql.DB

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
	log.Println("OnMessage", mcmd, scmd, string(data))

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
	}

	return define.ErrUnknownSubCmd
}

// ChangeUserTreasure 改变用户财富
func (p *Processor) ChangeUserTreasure(id int, score int64, varScore int64, diamond int64, varDiamond int64, changeType int) error {
	// 当前分数钻石
	if score < 0 || diamond < 0 {
		if err := GAME.QueryRow("SELECT user_score, user_diamond FROM user_treasure WHERE user_id = ?", id).Scan(&score, &diamond); err != nil {
			return err
		}
	}

	// 更新分数钻石
	if _, err := GAME.Exec("UPDATE user_treasure SET user_score = user_score + ?, user_diamond = user_diamond + ? WHERE user_id = ?", varScore, varDiamond, id); err != nil {
		return err
	}

	// 记录财富日志
	if _, err := LOG.Exec(fmt.Sprintf("INSERT INTO user_treasure_log_%d (user_id, cur_score, var_score, cur_diamond, var_diamond, change_type) VALUES (?, ?, ?, ?, ?, ?)", utils.Date()), id, score, varScore, diamond, varDiamond, changeType); err != nil {
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
	if err := GAME.QueryRow("SELECT user_id, user_level, bind_phone, user_score, user_diamond FROM view_information_treasure WHERE user_account = ?", fastRegister.Account).Scan(
		&replyFastRegister.UserID,
		&replyFastRegister.UserLevel,
		&replyFastRegister.BindPhone,
		&replyFastRegister.UserScore,
		&replyFastRegister.UserDiamond,
	); err == sql.ErrNoRows {
		// 插入用户信息
		res, err := GAME.Exec("INSERT INTO user_information (user_account, user_name, user_icon, user_gender, register_ip, register_machine) VALUES (?, ?, ?, ?, ?, ?)",
			fastRegister.Account,
			fastRegister.Name,
			fastRegister.Icon,
			fastRegister.Gender,
			fastRegister.IP,
			fastRegister.Machine,
		)
		if err != nil {
			return err
		}

		// 获取用户编号
		uid, err := res.LastInsertId()
		if err != nil {
			return err
		}

		replyFastRegister.UserID = int(uid)

		// 插入用户财富
		if _, err = GAME.Exec("INSERT INTO user_treasure (user_id) VALUES (?)", uid); err != nil {
			return err
		}

		// 用户初始分数钻石
		var score, diamond int64

		if err := GAME.QueryRow(`SELECT Content FROM game_config WHERE Title = "InitScore"`).Scan(&score); err != nil {
			return err
		}

		if err := GAME.QueryRow(`SELECT Content FROM game_config WHERE Title = "InitDiamond"`).Scan(&diamond); err != nil {
			return err
		}

		replyFastRegister.UserScore = score
		replyFastRegister.UserDiamond = diamond

		// 用户财富变化
		if err := p.ChangeUserTreasure(int(uid), 0, score, 0, diamond, define.ChangeTypeRegister); err != nil {
			return err
		}

		// 初始用户等级
		replyFastRegister.UserLevel = 1
	} else if err != nil {
		return err
	}

	log.Println(fastRegister, replyFastRegister)

	return replyFastRegister
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
func NewProcessor(server *network.Server) *Processor {
	var err error

	// todo SetMaxOpenConns, SetMaxIdleConns

	if LOG, err = sql.Open("mysql", "root:@/log"); err != nil {
		log.Println("Open log", err)
		return nil
	}

	if err = LOG.Ping(); err != nil {
		log.Println("Ping log", err)
		return nil
	}

	if GAME, err = sql.Open("mysql", "root:@/game"); err != nil {
		log.Println("Open game", err)
		return nil
	}

	if err = GAME.Ping(); err != nil {
		log.Println("Ping game", err)
		return nil
	}

	return &Processor{
		server: server,
	}
}
