package define

const (
	// ServiceProxy 代理服务
	ServiceProxy = 1

	// ServiceLogin 登陆服务
	ServiceLogin = 2

	// ServiceGame 游戏服务
	ServiceGame = 3
)

const (
	// ManagerCommon 通用
	ManagerCommon = 1
)

const (
	// ManagerRegisterService 注册服务
	ManagerRegisterService = 1
)

// RegisterService 注册服务
type RegisterService struct {
	ID          int    // 编号
	IP          string // 地址
	Count       int    // 计数
	GameType    int    // 游戏类型
	GameLevel   int    // 游戏等级
	ServiceType int    // 服务类型
	IsServe     bool   // 是否服务
}
