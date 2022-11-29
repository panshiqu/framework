package define

import (
	"net"
)

// Token 令牌
const Token = "13526535277"

// LengthLimit 长度限制
const LengthLimit = 8192

const (
	// ServiceProxy 代理服务
	ServiceProxy = 1

	// ServiceLogin 登陆服务
	ServiceLogin = 2

	// ServiceGame 游戏服务
	ServiceGame = 3
)

const (
	// GameUnknown 未知
	GameUnknown = 0

	// GameLandlords 斗地主
	GameLandlords = 1

	// GameFiveInARow 五子棋
	GameFiveInARow = 2
)

const (
	// LevelUnknown 未知
	LevelUnknown = 0

	// LevelOne 新手场
	LevelOne = 1

	// LevelTwo 初级场
	LevelTwo = 2

	// LevelThree 中级场
	LevelThree = 3

	// LevelFour 高级场
	LevelFour = 4
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
	// GlobalCommon 通用
	GlobalCommon = 0
)

const (
	// GlobalKeepAlive 保活
	GlobalKeepAlive = 0
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

	// ManagerNotifyChangeService 通知改变服务
	ManagerNotifyChangeService = 104

	// ManagerNotifyAllService 通知所有服务
	ManagerNotifyAllService = 105

	// ManagerNotifyIncrService 通知增加服务
	ManagerNotifyIncrService = 106

	// ManagerNotifyDecrService 通知删除服务
	ManagerNotifyDecrService = 107
)

const (
	// LoginCommon 通用
	LoginCommon = 1
)

const (
	// LoginFastRegister 快速注册
	LoginFastRegister = 1

	// LoginSignInDays 签到天数
	LoginSignInDays = 2

	// LoginSignIn 签到
	LoginSignIn = 3
)

const (
	// GameCommon 通用
	GameCommon = 100

	// GameTable 桌子
	GameTable = 150
)

const (
	// GameFastLogin 快速登陆
	GameFastLogin = 1

	// GameLogout 登出
	GameLogout = 2

	// GameReady 准备
	GameReady = 3

	// GameNotifySitDown 通知坐下
	GameNotifySitDown = 101

	// GameNotifyStandUp 通知站起
	GameNotifyStandUp = 102

	// GameNotifyStatus 通知状态
	GameNotifyStatus = 103

	// GameNotifyTreasure 通知财富
	GameNotifyTreasure = 104
)

const (
	// DBCommon 通用
	DBCommon = 1
)

const (
	// DBFastRegister 快速注册
	DBFastRegister = 1

	// DBFastLogin 快速登陆
	DBFastLogin = 2

	// DBChangeTreasure 改变财富
	DBChangeTreasure = 3

	// DBSignInDays 签到天数
	DBSignInDays = 4

	// DBSignIn 签到
	DBSignIn = 5

	// DBInsertOnlineCache 插入在线缓存
	DBInsertOnlineCache = 6

	// DBDeleteOnlineCache 删除在线缓存
	DBDeleteOnlineCache = 7

	// DBClearOnlineCache 清空在线缓存
	DBClearOnlineCache = 8
)

const (
	// GenderMale 男
	GenderMale = 1

	// GenderFemale 女
	GenderFemale = 2

	// GenderUnknown 未知
	GenderUnknown = 3
)

const (
	// ChangeTypeRegister 注册
	ChangeTypeRegister = 1

	// ChangeTypeWinLose 输赢
	ChangeTypeWinLose = 2

	// ChangeTypeSignIn 签到
	ChangeTypeSignIn = 3
)

const (
	// KeepAliveDead 死亡
	KeepAliveDead = 1

	// KeepAliveWarn 警告
	KeepAliveWarn = 2

	// KeepAliveSafe 安全
	KeepAliveSafe = 3
)

const (
	// TableStatusFree 空闲
	TableStatusFree = 0

	// TableStatusGame 游戏
	TableStatusGame = 1
)

const (
	// UserStatusFree 空闲
	UserStatusFree = 0

	// UserStatusReady 准备
	UserStatusReady = 1

	// UserStatusPlaying 游戏
	UserStatusPlaying = 2

	// UserStatusOffline 离线
	UserStatusOffline = 3
)

// InvalidChair 无效椅子
const InvalidChair = -1

// InvalidTable 无效桌子
const InvalidTable = -1

// TimerPerUser 用户持有
const TimerPerUser = 100

// TimerPerTable 桌子持有
const TimerPerTable = 100000

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

// GameService 游戏服务
type GameService struct {
	GameID    int `json:",omitempty"` // 编号
	GameType  int `json:",omitempty"` // 类型
	GameLevel int `json:",omitempty"` // 等级
}

// UserInfo 用户信息
type UserInfo struct {
	UserID      int    `json:",omitempty"` // 编号
	UserName    string `json:",omitempty"` // 名称
	UserIcon    int    `json:",omitempty"` // 图标
	UserLevel   int    `json:",omitempty"` // 等级
	UserGender  int    `json:",omitempty"` // 性别
	BindPhone   string `json:",omitempty"` // 绑定手机
	UserScore   int64  `json:",omitempty"` // 分数
	UserDiamond int64  `json:",omitempty"` // 钻石
}

// LoginCache 登陆缓存
type LoginCache struct {
	UserID int `json:",omitempty"` // 编号
}

// NotifySitDown 通知坐下
type NotifySitDown struct {
	UserInfo
	TableID    int `json:",omitempty"` // 桌子编号
	ChairID    int `json:",omitempty"` // 椅子编号
	UserStatus int `json:",omitempty"` // 用户状态
}

// NotifyStandUp 通知站起
type NotifyStandUp struct {
	ChairID int `json:",omitempty"` // 椅子编号
}

// NotifyStatus 通知状态
type NotifyStatus struct {
	ChairID    int `json:",omitempty"` // 椅子编号
	UserStatus int `json:",omitempty"` // 用户状态
}

// NotifyTreasure 通知财富
type NotifyTreasure struct {
	UserID     int   `json:",omitempty"` // 编号
	VarScore   int64 `json:",omitempty"` // 分数
	VarDiamond int64 `json:",omitempty"` // 钻石
	ChangeType int   `json:",omitempty"` // 类型
}

// FastRegister 快速注册
type FastRegister struct {
	Account  string `json:",omitempty"` // 账户
	Password string `json:",omitempty"` // 密码（未使用）
	Machine  string `json:",omitempty"` // 机器码
	Name     string `json:",omitempty"` // 名称
	Icon     int    `json:",omitempty"` // 图标
	Gender   int    `json:",omitempty"` // 性别
	IP       string `json:",omitempty"` // 地址
}

// ReplyFastRegister 回复快速注册
type ReplyFastRegister struct {
	ErrField
	UserInfo
	GameService
}

// ReplySignInDays 回复签到天数
type ReplySignInDays struct {
	ErrField
	Can  bool `json:",omitempty"` // 是否可以签到
	Days int  `json:",omitempty"` // 天数
}

// ReplySignIn 回复签到
type ReplySignIn struct {
	ErrField
	TotalDays     int   `json:",omitempty"` // 总天数
	ScoreReward   int64 `json:",omitempty"` // 分数奖励
	DiamondReward int64 `json:",omitempty"` // 钻石奖励
}

// FastLogin 快速登陆
type FastLogin struct {
	GameService
	UserID    int    `json:",omitempty"` // 编号
	Timestamp int64  `json:",omitempty"` // 时间戳
	Signature string `json:",omitempty"` // 加密签名
}

// ReplyFastLogin 回复快速登陆
type ReplyFastLogin struct {
	ErrField
	UserInfo
	GameService
	IsRobot bool `json:",omitempty"` // 机器人
}

// OnlineCache 在线缓存
type OnlineCache struct {
	GameService
	UserID int `json:",omitempty"` // 编号
}
