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

	// ManagerUpdateCount 更新计数
	ManagerUpdateCount = 2

	// ManagerOpenService 开启服务
	ManagerOpenService = 3

	// ManagerShutService 关闭服务
	ManagerShutService = 4

	// ManagerNotifyCurService 通知当前服务
	ManagerNotifyCurService = 101

	// ManagerNotifyAddService 通知增加服务
	ManagerNotifyAddService = 102

	// ManagerNotifyDelService 通知删除服务
	ManagerNotifyDelService = 103
)

const (
	// LoginCommon 通用
	LoginCommon = 1
)

const (
	// LoginFastRegister 快速注册
	LoginFastRegister = 1
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

// FastRegister 快速注册
type FastRegister struct {
	Account  string `json:",omitempty"` // 账户
	Password string `json:",omitempty"` // 密码
	Machine  string `json:",omitempty"` // 机器码
	Name     string `json:",omitempty"` // 名称
	Icon     int    `json:",omitempty"` // 图标
	Gender   int    `json:",omitempty"` // 性别
}
