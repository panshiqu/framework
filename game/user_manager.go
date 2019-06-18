package game

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"../define"
)

var uins UserManager

// UserManager 用户管理
type UserManager struct {
	mutex sync.Mutex
	users map[uint32]*UserItem
}

// Delete 删除用户
func (u *UserManager) Delete(id uint32) {
	u.mutex.Lock()
	if userItem, ok := u.users[id]; ok {
		if err := userItem.WriteToDB(userItem.CacheScore(), userItem.CacheDiamond(), define.ChangeTypeWinLose); err != nil {
			log.Println("UserManager WriteToDB", err)
		}
	}
	delete(u.users, id)
	u.mutex.Unlock()
}

// Search 查找用户
func (u *UserManager) Search(id uint32) *UserItem {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if userItem, ok := u.users[id]; ok {
		return userItem
	}

	return nil
}

// Insert 插入用户
func (u *UserManager) Insert(conn net.Conn, reply *define.ReplyFastLogin) *UserItem {
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

// Monitor 监视器
func (u *UserManager) Monitor(w http.ResponseWriter, r *http.Request) {
	u.mutex.Lock()
	fmt.Fprintln(w, "online:", len(u.users))
	u.mutex.Unlock()
}
