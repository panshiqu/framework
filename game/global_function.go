package game

import (
	"log"

	"github.com/panshiqu/framework/define"
)

// UpdateOnlineCache 更新在线缓存，插入、删除、清空
func UpdateOnlineCache(ns ...int) {
	cache := &define.OnlineCache{
		GameID: define.CG.ID,
	}

	scmd := define.DBClearOnlineCache
	if len(ns) == 2 {
		scmd = ns[0]
		cache.UserID = ns[1]
	}

	if scmd == define.DBInsertOnlineCache {
		cache.GameType = define.CG.GameType
		cache.GameLevel = define.CG.GameLevel
	}

	if err := rpc.JSONCall(define.DBCommon, uint16(scmd), cache, nil); err != nil {
		log.Println("OnlineCache", scmd, cache, err)
	}
}

// ClearOnlineCacheNowAndExit 存在只是为了精简入口函数调用
func ClearOnlineCacheNowAndExit() func(...int) {
	UpdateOnlineCache()
	return UpdateOnlineCache
}
