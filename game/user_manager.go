package game

import (
	"net"
	"sync"

	"github.com/panshiqu/framework/define"
)

// UserManager 用户管理
type UserManager struct {
	mutex sync.Mutex
	users map[int]*UserItem
}

// Remove 移除用户
func (u *UserManager) Remove(id int) {
	u.mutex.Lock()
	delete(u.users, id)
	u.mutex.Unlock()
}

// Create 创造用户
func (u *UserManager) Create(conn net.Conn, reply *define.ReplyFastLogin) *UserItem {
	userItem := &UserItem{
		id:      reply.UserID,
		name:    reply.UserName,
		icon:    reply.UserIcon,
		level:   reply.UserLevel,
		gender:  reply.UserGender,
		score:   reply.UserScore,
		diamond: reply.UserDiamond,
		robot:   reply.IsRobot,
		conn:    conn,
	}

	u.mutex.Lock()
	u.users[userItem.UserID()] = userItem
	u.mutex.Unlock()

	return userItem
}
