package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/network"
	"github.com/panshiqu/framework/utils"
)

// LOG 日志数据库
var LOG *sql.DB

// GAME 游戏数据库
var GAME *sql.DB

// REDIS 缓存连接池
var REDIS *redis.Pool

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
		return utils.Wrap(err)
	}

	// 实现快捷回复消息
	return utils.Wrap(network.SendJSONMessage(conn, mcmd, scmd, ret))
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
	case define.DBFastLogin:
		return p.OnSubFastLogin(conn, data)
	case define.DBChangeTreasure:
		return p.OnSubChangeTreasure(conn, data)
	case define.DBSignInDays:
		return p.OnSubSignInDays(conn, data)
	case define.DBSignIn:
		return p.OnSubSignIn(conn, data)
	}

	return define.ErrUnknownSubCmd
}

// ChangeUserTreasure 改变用户财富
func (p *Processor) ChangeUserTreasure(id int, score int64, varScore int64, diamond int64, varDiamond int64, changeType int) error {
	// 当前分数钻石
	if score < 0 || diamond < 0 {
		if err := GAME.QueryRow("SELECT user_score, user_diamond FROM user_treasure WHERE user_id = ?", id).Scan(&score, &diamond); err != nil {
			return utils.Wrap(err)
		}
	}

	tx, err := GAME.Begin()
	if err != nil {
		return utils.Wrap(err)
	}
	defer tx.Rollback()

	// 更新分数钻石
	if _, err := tx.Exec("UPDATE user_treasure SET user_score = user_score + ?, user_diamond = user_diamond + ? WHERE user_id = ?", varScore, varDiamond, id); err != nil {
		return utils.Wrap(err)
	}

	// 记录财富日志
	if _, err := tx.Exec(fmt.Sprintf("INSERT INTO log.user_treasure_log_%s (user_id, cur_score, var_score, cur_diamond, var_diamond, change_type) VALUES (?, ?, ?, ?, ?, ?)", time.Now().Format("20060102")), id, score, varScore, diamond, varDiamond, changeType); err != nil {
		return utils.Wrap(err)
	}

	if err := tx.Commit(); err != nil {
		return utils.Wrap(err)
	}

	return nil
}

// OnSubFastRegister 快速注册子命令
func (p *Processor) OnSubFastRegister(conn net.Conn, data []byte) interface{} {
	fastRegister := &define.FastRegister{}
	replyFastRegister := &define.ReplyFastRegister{}

	if err := json.Unmarshal(data, fastRegister); err != nil {
		return utils.Wrap(err)
	}

	// 查询用户信息
	if err := GAME.QueryRow("SELECT user_id, user_level, bind_phone, user_score, user_diamond FROM view_information_treasure WHERE user_account = ?", fastRegister.Account).Scan(
		&replyFastRegister.UserID,
		&replyFastRegister.UserLevel,
		&replyFastRegister.BindPhone,
		&replyFastRegister.UserScore,
		&replyFastRegister.UserDiamond,
	); errors.Is(err, sql.ErrNoRows) {
		tx, err := GAME.Begin()
		if err != nil {
			return utils.Wrap(err)
		}
		defer tx.Rollback()

		// 插入用户信息
		res, err := tx.Exec("INSERT INTO user_information (user_account, user_name, user_icon, user_gender, register_ip, register_machine) VALUES (?, ?, ?, ?, ?, ?)",
			fastRegister.Account,
			fastRegister.Name,
			fastRegister.Icon,
			fastRegister.Gender,
			fastRegister.IP,
			fastRegister.Machine,
		)
		if err != nil {
			return utils.Wrap(err)
		}

		// 获取用户编号
		uid, err := res.LastInsertId()
		if err != nil {
			return utils.Wrap(err)
		}

		replyFastRegister.UserID = int(uid)

		// 插入用户财富
		if _, err = tx.Exec("INSERT INTO user_treasure (user_id) VALUES (?)", uid); err != nil {
			return utils.Wrap(err)
		}

		// 插入签到记录
		if _, err = tx.Exec("INSERT INTO sign_in_record (user_id) VALUES (?)", uid); err != nil {
			return utils.Wrap(err)
		}

		if err := tx.Commit(); err != nil {
			return utils.Wrap(err)
		}

		// 用户初始分数钻石
		var score, diamond int64

		if err := GAME.QueryRow(`SELECT Content FROM game_config WHERE Title = "InitScore"`).Scan(&score); err != nil {
			return utils.Wrap(err)
		}

		if err := GAME.QueryRow(`SELECT Content FROM game_config WHERE Title = "InitDiamond"`).Scan(&diamond); err != nil {
			return utils.Wrap(err)
		}

		replyFastRegister.UserScore = score
		replyFastRegister.UserDiamond = diamond

		// 用户财富变化
		if err := p.ChangeUserTreasure(int(uid), 0, score, 0, diamond, define.ChangeTypeRegister); err != nil {
			return utils.Wrap(err)
		}

		// 初始用户等级
		replyFastRegister.UserLevel = 1
	} else if err != nil {
		return utils.Wrap(err)
	}

	// 总是更新这些字段
	if _, err := GAME.Exec("UPDATE user_information SET user_name = ?, user_icon = ?, user_gender = ? WHERE user_id = ?",
		fastRegister.Name,
		fastRegister.Icon,
		fastRegister.Gender,
		replyFastRegister.UserID,
	); err != nil {
		return utils.Wrap(err)
	}

	return replyFastRegister
}

// OnSubSignInDays 签到天数
func (p *Processor) OnSubSignInDays(conn net.Conn, data []byte) interface{} {
	replySignInDays := &define.ReplySignInDays{}

	if err := GAME.QueryRow("CALL procedure_user_sign_in_days(?)", data).Scan(&replySignInDays.Can, &replySignInDays.Days); err != nil {
		return utils.Wrap(err)
	}

	return replySignInDays
}

// OnSubSignIn 签到
func (p *Processor) OnSubSignIn(conn net.Conn, data []byte) interface{} {
	replySignIn := &define.ReplySignIn{}

	if err := GAME.QueryRow("CALL procedure_user_sign_in(?)", data).Scan(&replySignIn.Errno, &replySignIn.Errdesc,
		&replySignIn.TotalDays, &replySignIn.ScoreReward, &replySignIn.DiamondReward); err != nil {
		return utils.Wrap(err)
	}

	return replySignIn
}

// OnSubFastLogin 快速登陆子命令
func (p *Processor) OnSubFastLogin(conn net.Conn, data []byte) interface{} {
	replyFastLogin := &define.ReplyFastLogin{}

	// 查询用户信息
	if err := GAME.QueryRow("SELECT user_id, user_name, user_icon, user_level, user_gender+0, bind_phone, user_score, user_diamond, is_robot FROM view_information_treasure WHERE user_id = ?", data).Scan(
		&replyFastLogin.UserID,
		&replyFastLogin.UserName,
		&replyFastLogin.UserIcon,
		&replyFastLogin.UserLevel,
		&replyFastLogin.UserGender,
		&replyFastLogin.BindPhone,
		&replyFastLogin.UserScore,
		&replyFastLogin.UserDiamond,
		&replyFastLogin.IsRobot,
	); err != nil {
		return utils.Wrap(err)
	}

	return replyFastLogin
}

// OnSubChangeTreasure 改变财富
func (p *Processor) OnSubChangeTreasure(conn net.Conn, data []byte) interface{} {
	notifyTreasure := &define.NotifyTreasure{}

	if err := json.Unmarshal(data, notifyTreasure); err != nil {
		return utils.Wrap(err)
	}

	// 用户财富变化
	return utils.Wrap(p.ChangeUserTreasure(notifyTreasure.UserID,
		-1, notifyTreasure.VarScore,
		-1, notifyTreasure.VarDiamond,
		notifyTreasure.ChangeType))
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
func NewProcessor(server *network.Server, config *define.ConfigDB) *Processor {
	var err error

	// todo SetMaxOpenConns, SetMaxIdleConns

	if LOG, err = sql.Open("mysql", config.LogDSN); err != nil {
		log.Println("Open log", err)
		return nil
	}

	if err = LOG.Ping(); err != nil {
		log.Println("Ping log", err)
		return nil
	}

	if GAME, err = sql.Open("mysql", config.GameDSN); err != nil {
		log.Println("Open game", err)
		return nil
	}

	if err = GAME.Ping(); err != nil {
		log.Println("Ping game", err)
		return nil
	}

	REDIS = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: time.Hour,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(config.RedisURL)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	return &Processor{
		server: server,
	}
}

func GetRedis(database int) (conn redis.Conn) {
	for {
		conn = REDIS.Get()
		if _, err := conn.Do("SELECT", database); err != nil {
			log.Println("Redis select", err)
			conn.Close()
			continue
		}
		break
	}

	prefix := "Unknown"
	pc := make([]uintptr, 1)
	if n := runtime.Callers(2, pc); n == 1 {
		// the Name method can be called with nil
		prefix = path.Base(runtime.FuncForPC(pc[0]).Name())
	}

	return redis.NewLoggingConn(conn, log.Default(), prefix)
}
