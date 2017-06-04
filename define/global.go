package define

import (
	"net"
)

const (
	// ServiceProxy 代理服务
	ServiceProxy = 1

	// ServiceLogin 登陆服务
	ServiceLogin = 2

	// ServiceGame 游戏服务
	ServiceGame = 3
)

const (
	// CapacityProxy 代理容量
	CapacityProxy = 1000

	// CapacityLogin 登陆容量
	CapacityLogin = 1000

	// CapacityGame 游戏容量
	CapacityGame = 1000
)

const (
	// ManagerCommon 通用
	ManagerCommon = 1
)

const (
	// ManagerRegisterService 注册服务
	ManagerRegisterService = 1

	// ManagerUpdateServiceCount 更新服务计数
	ManagerUpdateServiceCount = 2

	// ManagerOpenService 开启服务
	ManagerOpenService = 3

	// ManagerShutService 关闭服务
	ManagerShutService = 4
)

// Service 服务
type Service struct {
	ID          int      `json:",omitempty"` // 编号
	IP          string   `json:",omitempty"` // 地址
	Count       int      `json:",omitempty"` // 计数
	GameType    int      `json:",omitempty"` // 游戏类型
	GameLevel   int      `json:",omitempty"` // 游戏等级
	ServiceType int      `json:",omitempty"` // 服务类型
	IsServe     bool     `json:",omitempty"` // 是否服务
	Conn        net.Conn `json:"-"`          // 网络连接
}

// ServiceCount 服务计数
type ServiceCount struct {
	ID    int `json:",omitempty"` // 编号
	Count int `json:",omitempty"` // 计数
}

// OpenService 开启服务
type OpenService struct {
	ID int `json:",omitempty"` // 编号
}
