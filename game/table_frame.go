package game

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"../define"
)

// TableFrame 桌子框架
type TableFrame struct {
	id     int         // 编号
	status int32       // 状态
	mutex  sync.Mutex  // 加锁
	users  []*UserItem // 用户

	table define.ITableLogic // 桌子逻辑
}

// TableID 桌子编号
func (t *TableFrame) TableID() int {
	return t.id
}

// TableStatus 桌子状态
func (t *TableFrame) TableStatus() int32 {
	return atomic.LoadInt32(&t.status)
}

// SetTableLogic 设置桌子逻辑
func (t *TableFrame) SetTableLogic(v define.ITableLogic) {
	t.table = v
}

// TableUser 桌子用户
func (t *TableFrame) TableUser(chair int) *UserItem {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if chair < define.CG.UserPerTable {
		return t.users[chair]
	}
	return nil
}

// GetUser 获取用户
func (t *TableFrame) GetUser(chair int) define.IUserItem {
	return t.TableUser(chair)
}

// UserCount 用户数量
func (t *TableFrame) UserCount() (cnt int) {
	for i := 0; i < define.CG.UserPerTable; i++ {
		if t.TableUser(i) != nil {
			cnt++
		}
	}

	return
}

// ReadyCount 准备数量
func (t *TableFrame) ReadyCount() (cnt int) {
	for i := 0; i < define.CG.UserPerTable; i++ {
		if user := t.TableUser(i); user != nil && user.UserStatus() == define.UserStatusReady {
			cnt++
		}
	}

	return
}

// SetUserStatus 设置用户状态
func (t *TableFrame) SetUserStatus(status int) {
	for i := 0; i < define.CG.UserPerTable; i++ {
		if user := t.TableUser(i); user != nil {
			user.SetUserStatus(status)
		}
	}
}

// SendTableUserInfo 发送同桌玩家信息
func (t *TableFrame) SendTableUserInfo(userItem *UserItem) {
	for i := 0; i < define.CG.UserPerTable; i++ {
		if i == userItem.ChairID() {
			continue
		}

		if user := t.TableUser(i); user != nil {
			userItem.SendJSONMessage(define.GameCommon, define.GameNotifySitDown, user.TableUserInfo())
		}
	}
}

// SitDown 坐下
func (t *TableFrame) SitDown(userItem *UserItem) {
	t.mutex.Lock()
	chair := define.InvalidChair
	for k, v := range t.users {
		if v == nil {
			chair = k
			break
		}
	}
	t.users[chair] = userItem
	userItem.SetChairID(chair)
	userItem.SetTableFrame(t)
	t.mutex.Unlock()

	// 广播我的坐下
	t.SendTableJSONMessage(define.GameCommon, define.GameNotifySitDown, userItem.TableUserInfo())

	// 发送同桌玩家信息
	t.SendTableUserInfo(userItem)

	// 用户坐下
	if err := t.table.OnUserSitDown(userItem); err != nil {
		log.Println("TableFrame OnUserSitDown", err)
	}
}

// StandUp 站起
func (t *TableFrame) StandUp(userItem *UserItem) {
	// 用户站起
	if err := t.table.OnUserStandUp(userItem); err != nil {
		log.Println("TableFrame OnUserStandUp", err)
	}

	t.mutex.Lock()
	chair := userItem.ChairID()
	t.users[chair] = nil
	userItem.SetChairID(define.InvalidChair)
	userItem.SetTableFrame(nil)
	t.mutex.Unlock()

	standUp := &define.NotifyStandUp{
		ChairID: chair,
	}

	// 广播用户站起
	t.SendTableJSONMessage(define.GameCommon, define.GameNotifyStandUp, standUp)
}

// Reconnect 重连
func (t *TableFrame) Reconnect(userItem *UserItem) {
	if err := t.table.OnUserReconnect(userItem); err != nil {
		log.Println("TableFrame OnUserReconnect", err)
	}
}

// StartGame 开始游戏
func (t *TableFrame) StartGame() {
	// 检查准备数量
	if t.ReadyCount() < define.CG.MinReadyStart {
		return
	}

	// 设置桌子状态
	if !atomic.CompareAndSwapInt32(&t.status, define.TableStatusFree, define.TableStatusGame) {
		return
	}

	// 设置游戏状态
	t.SetUserStatus(define.UserStatusPlaying)

	// 游戏开始
	if err := t.table.OnGameStart(); err != nil {
		log.Println("TableFrame OnGameStart", err)
	}
}

// ConcludeGame 结束游戏
func (t *TableFrame) ConcludeGame() {
	// 设置桌子状态
	if !atomic.CompareAndSwapInt32(&t.status, define.TableStatusGame, define.TableStatusFree) {
		return
	}

	for i := 0; i < define.CG.UserPerTable; i++ {
		if user := t.TableUser(i); user != nil && user.UserStatus() == define.UserStatusOffline {
			// 用户站起
			t.StandUp(user)

			// 删除用户
			uins.Delete(user.UserID())
		}
	}

	// 设置空闲状态
	t.SetUserStatus(define.UserStatusFree)
}

// AddTimer 添加定时器
func (t *TableFrame) AddTimer(id int, duration time.Duration, parameter interface{}, persistence bool) {
	if id >= 0 && id < define.TimerPerUser {
		sins.Add(t.TableID()*define.TimerPerTable+define.TimerPerTable+id, duration, parameter, persistence)
	}
}

// RunAfter 添加定时器
func (t *TableFrame) RunAfter(id int, duration time.Duration, parameter interface{}) {
	t.AddTimer(id, duration, parameter, false)
}

// RunAlways 添加定时器
func (t *TableFrame) RunAlways(id int, duration time.Duration, parameter interface{}) {
	t.AddTimer(id, duration, parameter, true)
}

// RemoveTimer 移除定时器
func (t *TableFrame) RemoveTimer(id int) {
	if id >= 0 && id < define.TimerPerUser {
		sins.Remove(t.TableID()*define.TimerPerTable + define.TimerPerTable + id)
	}
}

// SurplusDuration 定时器剩余时间
func (t *TableFrame) SurplusDuration(id int) time.Duration {
	if id >= 0 && id < define.TimerPerUser {
		return sins.Surplus(t.TableID()*define.TimerPerTable + define.TimerPerTable + id)
	}

	return 0
}

// OnTimer 定时器
func (t *TableFrame) OnTimer(id int, parameter interface{}) error {
	if id < define.TimerPerUser {
		return t.table.OnTimer(id, parameter)
	}

	if user := t.TableUser((id - define.TimerPerUser) / define.TimerPerUser); user != nil {
		return user.OnTimer((id-define.TimerPerUser)%define.TimerPerUser, parameter)
	}

	return nil
}

// OnMessage 收到消息
func (t *TableFrame) OnMessage(scmd uint16, data []byte, userItem *UserItem) error {
	return t.table.OnMessage(scmd, data, userItem)
}

// SendTableMessage 发送桌子消息
func (t *TableFrame) SendTableMessage(mcmd uint16, scmd uint16, data []byte) {
	for i := 0; i < define.CG.UserPerTable; i++ {
		t.SendChairMessage(i, mcmd, scmd, data)
	}
}

// SendTableJSONMessage 发送桌子消息
func (t *TableFrame) SendTableJSONMessage(mcmd uint16, scmd uint16, js interface{}) {
	if data, err := json.Marshal(js); err == nil {
		t.SendTableMessage(mcmd, scmd, data)
	}
}

// SendChairMessage 发送椅子消息
func (t *TableFrame) SendChairMessage(chair int, mcmd uint16, scmd uint16, data []byte) {
	if user := t.TableUser(chair); user != nil {
		user.SendMessage(mcmd, scmd, data)
	}
}

// SendChairJSONMessage 发送椅子消息
func (t *TableFrame) SendChairJSONMessage(chair int, mcmd uint16, scmd uint16, js interface{}) {
	if data, err := json.Marshal(js); err == nil {
		t.SendChairMessage(chair, mcmd, scmd, data)
	}
}

// Monitor 监视器
func (t *TableFrame) Monitor(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "id:%3d, status:%d, usercount:%3d\n", t.id, t.TableStatus(), t.UserCount())
	for i := 0; i < define.CG.UserPerTable; i++ {
		if user := t.TableUser(i); user != nil {
			fmt.Fprintf(w, "\tid:%8d, score:%10d, diamond:%8d, status:%d, chair:%3d, name:%s\n", user.UserID(), user.UserScore(), user.UserDiamond(), user.UserStatus(), user.ChairID(), user.UserName())
		}
	}
}
