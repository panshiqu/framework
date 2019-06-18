package define

// IUserItem 用户接口
type IUserItem interface {
	// 用户编号
	UserID() uint32

	// 用户分数
	UserScore() int64

	// 用户钻石
	UserDiamond() int64

	// 是否机器人
	IsRobot() bool

	// 椅子编号
	ChairID() int

	// 写入财富
	WriteTreasure(int64, int64, int) error

	// 发送消息
	SendMessage(uint16, uint16, []byte)

	// 发送消息
	SendJSONMessage(uint16, uint16, interface{})
}

// ITableFrame 桌子框架接口
type ITableFrame interface {
	// 桌子编号
	TableID() int

	// 获取用户
	GetUser(int) IUserItem

	// 结束游戏
	ConcludeGame()

	// 发送桌子消息
	SendTableMessage(uint16, uint16, []byte)

	// 发送桌子消息
	SendTableJSONMessage(uint16, uint16, interface{})

	// 发送椅子消息
	SendChairMessage(int, uint16, uint16, []byte)

	// 发送椅子消息
	SendChairJSONMessage(int, uint16, uint16, interface{})
}

// ITableLogic 桌子逻辑接口
type ITableLogic interface {
	// 初始化
	OnInit() error

	// 游戏开始
	OnGameStart() error

	// 游戏结束
	OnGameConclude() error

	// 用户坐下
	OnUserSitDown(IUserItem) error

	// 用户站起
	OnUserStandUp(IUserItem) error

	// 用户重连
	OnUserReconnect(IUserItem) error

	// 收到消息
	OnMessage(uint16, []byte, IUserItem) error

	// 定时器
	OnTimer(int, interface{}) error
}
