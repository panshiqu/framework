package game

import (
	"encoding/json"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"../define"
	"../network"
)

// UserItem 用户
type UserItem struct {
	id           int      // 编号
	name         string   // 名称
	icon         int      // 图标
	level        int      // 等级
	gender       uint8      // 性别
	phone        string   // 手机
	score        int64    // 分数
	cacheScore   int64    // 缓存分数
	diamond      int64    // 钻石
	cacheDiamond int64    // 缓存钻石
	robot        bool     // 机器人
	conn         net.Conn // 网络连接（暂不为它考虑加锁）

	status     int32       // 状态
	chairID    int         // 椅子编号
	tableFrame *TableFrame // 桌子框架

	mutex sync.Mutex // 财富加锁
}

// UserID 用户编号
func (u *UserItem) UserID() int {
	return u.id
}

// UserName 用户名称
func (u *UserItem) UserName() string {
	return u.name
}

// UserIcon 用户图标
func (u *UserItem) UserIcon() int {
	return u.icon
}

// UserLevel 用户等级
func (u *UserItem) UserLevel() int {
	return u.level
}

// UserGender 用户性别
func (u *UserItem) UserGender() uint8 {
	return u.gender
}

// BindPhone 用户手机
func (u *UserItem) BindPhone() string {
	return u.phone
}

// UserScore 用户分数
func (u *UserItem) UserScore() int64 {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	return u.score + u.cacheScore
}

// CacheScore 缓存分数
func (u *UserItem) CacheScore() int64 {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	return u.cacheScore
}

// UserDiamond 用户钻石
func (u *UserItem) UserDiamond() int64 {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	return u.diamond + u.cacheDiamond
}

// CacheDiamond 缓存钻石
func (u *UserItem) CacheDiamond() int64 {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	return u.cacheDiamond
}

// IsRobot 是否机器人
func (u *UserItem) IsRobot() bool {
	return u.robot
}

// SetConn 设置网络连接
func (u *UserItem) SetConn(v net.Conn) {
	u.conn = v
}

// UserStatus 用户状态
func (u *UserItem) UserStatus() int {
	return int(atomic.LoadInt32(&u.status))
}

// SetUserStatus 设置用户状态
func (u *UserItem) SetUserStatus(v int) {
	atomic.StoreInt32(&u.status, int32(v))

	if tableFrame := u.TableFrame(); tableFrame != nil {
		notifyStatus := &define.NotifyStatus{
			ChairID:    u.ChairID(),
			UserStatus: u.UserStatus(),
		}

		tableFrame.SendTableJSONMessage(define.GameCommon, define.GameNotifyStatus, notifyStatus)
	}
}

// ChairID 椅子编号
func (u *UserItem) ChairID() int {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	return u.chairID
}

// SetChairID 设置椅子编号
func (u *UserItem) SetChairID(v int) {
	u.mutex.Lock()
	u.chairID = v
	u.mutex.Unlock()
}

// TableID 桌子编号
func (u *UserItem) TableID() int {
	if tableFrame := u.TableFrame(); tableFrame != nil {
		return tableFrame.TableID()
	}

	return define.InvalidTable
}

// TableFrame 桌子框架
func (u *UserItem) TableFrame() *TableFrame {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	return u.tableFrame
}

// SetTableFrame 设置桌子框架
func (u *UserItem) SetTableFrame(v *TableFrame) {
	u.mutex.Lock()
	u.tableFrame = v
	u.mutex.Unlock()
}

// TableUserInfo 桌子用户信息
func (u *UserItem) TableUserInfo() *define.NotifySitDown {
	return &define.NotifySitDown{
		UserInfo: define.UserInfo{
			UserID:      u.id,
			UserName:    u.name,
			UserIcon:    u.icon,
			UserLevel:   u.level,
			UserGender:  uint8(u.gender),
			UserScore:   u.UserScore(),
			UserDiamond: u.UserDiamond(),
		},
		TableID:    u.TableID(),
		ChairID:    u.ChairID(),
		UserStatus: u.UserStatus(),
	}
}

// WriteScore 写入分数
func (u *UserItem) WriteScore(varScore int64, changeType int) error {
	return u.WriteTreasure(varScore, 0, changeType)
}

// WriteDiamond 写入钻石
func (u *UserItem) WriteDiamond(varDiamond int64, changeType int) error {
	return u.WriteTreasure(0, varDiamond, changeType)
}

// WriteTreasure 写入财富
func (u *UserItem) WriteTreasure(varScore int64, varDiamond int64, changeType int) (err error) {
	defer func() {
		if err == nil {
			if tableFrame := u.TableFrame(); tableFrame != nil {
				notifyTreasure := &define.NotifyTreasure{
					UserID:     u.id,
					VarScore:   varScore,
					VarDiamond: varDiamond,
					ChangeType: changeType,
				}

				tableFrame.SendTableJSONMessage(define.GameCommon, define.GameNotifyTreasure, notifyTreasure)
			}
		}
	}()

	u.mutex.Lock()
	defer u.mutex.Unlock()

	// 分数不足
	if u.score+u.cacheScore+varScore < 0 {
		return define.ErrNotEnoughScore
	}

	// 钻石不足
	if u.diamond+u.cacheDiamond+varDiamond < 0 {
		return define.ErrNotEnoughDiamond
	}

	// 缓存输赢
	if changeType == define.ChangeTypeWinLose {
		u.cacheScore += varScore
		u.cacheDiamond += varDiamond
		return nil
	}

	// 写入数据库
	if err := u.WriteToDB(varScore, varDiamond, changeType); err != nil {
		return err
	}

	// 更新财富
	u.score += varScore
	u.diamond += varDiamond

	return nil
}

// WriteToDB 写入数据库
func (u *UserItem) WriteToDB(varScore int64, varDiamond int64, changeType int) error {
	if varScore == 0 && varDiamond == 0 {
		return nil
	}

	notifyTreasure := &define.NotifyTreasure{
		UserID:     u.id,
		VarScore:   varScore,
		VarDiamond: varDiamond,
		ChangeType: changeType,
	}

	return rpc.JSONCall(define.DBCommon, define.DBChangeTreasure, notifyTreasure, nil)
}

// SendMessage 发送消息
func (u *UserItem) SendMessage(mcmd uint16, scmd uint16, data []byte) {
	if err := network.SendMessage(u.conn, mcmd, scmd, data); err != nil {
		log.Println("UserItem SendMessage", err)
	}
}

// SendJSONMessage 发送消息
func (u *UserItem) SendJSONMessage(mcmd uint16, scmd uint16, js interface{}) {
	if data, err := json.Marshal(js); err == nil {
		u.SendMessage(mcmd, scmd, data)
	}
}

// AddTimer 添加定时器
func (u *UserItem) AddTimer(id int, duration time.Duration, parameter interface{}, persistence bool) {
	if id >= 0 && id < define.TimerPerUser {
		sins.Add(u.TableID()*define.TimerPerTable+define.TimerPerTable+u.ChairID()*define.TimerPerUser+define.TimerPerUser+id, duration, parameter, persistence)
	}
}

// RunAfter 添加定时器
func (u *UserItem) RunAfter(id int, duration time.Duration, parameter interface{}) {
	u.AddTimer(id, duration, parameter, false)
}

// RunAlways 添加定时器
func (u *UserItem) RunAlways(id int, duration time.Duration, parameter interface{}) {
	u.AddTimer(id, duration, parameter, true)
}

// RemoveTimer 移除定时器
func (u *UserItem) RemoveTimer(id int) {
	if id >= 0 && id < define.TimerPerUser {
		sins.Remove(u.TableID()*define.TimerPerTable + define.TimerPerTable + u.ChairID()*define.TimerPerUser + define.TimerPerUser + id)
	}
}

// SurplusDuration 定时器剩余时间
func (u *UserItem) SurplusDuration(id int) time.Duration {
	if id >= 0 && id < define.TimerPerUser {
		return sins.Surplus(u.TableID()*define.TimerPerTable + define.TimerPerTable + u.ChairID()*define.TimerPerUser + define.TimerPerUser + id)
	}

	return 0
}

// OnTimer 定时器
func (u *UserItem) OnTimer(id int, parameter interface{}) error {
	return nil
}
