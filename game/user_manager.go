package game

import (
	"net"
	"sync"

	"github.com/panshiqu/framework/define"
)

var uins UserManager

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

// Search 查找用户
func (u *UserManager) Search(id int) *UserItem {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if userItem, ok := u.users[id]; ok {
		return userItem
	}

	return nil
}

// Create 创造用户
func (u *UserManager) Create(conn net.Conn, reply *define.ReplyFastLogin) *UserItem {
	userItem := &UserItem{
		id:      reply.UserID,
		name:    reply.UserName,
		icon:    reply.UserIcon,
		level:   reply.UserLevel,
		gender:  reply.UserGender,
		phone:   reply.BindPhone,
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
