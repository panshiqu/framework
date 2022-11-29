package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"path"
	"runtime"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/utils"
)

// LOG 日志数据库
var LOG *sql.DB

// GAME 游戏数据库
var GAME *sql.DB

// REDIS 缓存连接池
var REDIS *redis.Pool

func InitDatabase(logDSN, gameDSN, redisURL string) (err error) {
	if LOG, err = sql.Open("mysql", logDSN); err != nil {
		return utils.Wrap(err)
	}

	LOG.SetMaxIdleConns(100)
	LOG.SetConnMaxIdleTime(time.Hour)

	if err = LOG.Ping(); err != nil {
		return utils.Wrap(err)
	}

	if GAME, err = sql.Open("mysql", gameDSN); err != nil {
		return utils.Wrap(err)
	}

	GAME.SetMaxIdleConns(100)
	GAME.SetConnMaxIdleTime(time.Hour)

	if err = GAME.Ping(); err != nil {
		return utils.Wrap(err)
	}

	REDIS = &redis.Pool{
		MaxIdle:     100,
		IdleTimeout: time.Hour,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(redisURL)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

func GetRedis(database int) (conn redis.Conn) {
	// 也可照搬redis.errorConn错误推迟到下次调用
	for i := 5 * time.Millisecond; ; i *= 2 {
		conn = REDIS.Get()
		if _, err := conn.Do("SELECT", database); err != nil {
			log.Println("Redis select", err)
			if i > time.Second {
				i = time.Second
			}
			time.Sleep(i)
			conn.Close()
			continue
		}
		break
	}

	prefix := "Unknown"
	pc := make([]uintptr, 1)
	if n := runtime.Callers(2, pc); n == 1 {
		// the Name method can be called with nil
		prefix = path.Ext(runtime.FuncForPC(pc[0]).Name())
	}

	return redis.NewLoggingConn(conn, log.Default(), prefix)
}

// GetOnlineCache 获取在线缓存
func GetOnlineCache(id int, reply any) error {
	rc := GetRedis(define.RedisOnline)
	defer rc.Close()

	data, err := redis.Bytes(RedisGetKeys.Do(rc, fmt.Sprintf("Online_*_%d", id)))
	if errors.Is(err, redis.ErrNil) {
		return nil
	}
	if err != nil {
		return utils.Wrap(err)
	}

	return json.Unmarshal(data, reply)
}

// too many results to unpack
// Online_2_[0~9999], support data item 7778+, or replace keys with scan
var RedisDelKeys = redis.NewScript(1, `
local results = redis.call('keys', KEYS[1])
if #results == 0 then
	return 0
end
return redis.call('del', unpack(results))`)

var RedisGetKeys = redis.NewScript(1, `
local results = redis.call('keys', KEYS[1])
if #results == 0 then
	return nil
end
if #results > 1 then
	return redis.error_reply('more than one')
end
return redis.call('get', results[1])`)
