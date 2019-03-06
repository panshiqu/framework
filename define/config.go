package define

// ConfigManager 管理服务器配置
type ConfigManager struct {
	PprofIP  string
	ListenIP string
}

// ConfigProxy 代理服务器配置
type ConfigProxy struct {
	ID       int
	PprofIP  string
	ListenIP string
}

// ConfigLogin 登陆服务器配置
type ConfigLogin struct {
	ID       int
	DBIP     string
	DialIP   string
	PprofIP  string
	ListenIP string
}

// ConfigGame 游戏服务器配置
type ConfigGame struct {
	ID            int
	//DialIP        string
	PprofIP       string
	ListenIP      string
	GameType      int // 游戏类型
	UserPerTable  int // 用户数量每桌
	MinReadyStart int // 最小准备开始
}

// CG 游戏配置
var CG ConfigGame

// ConfigDB 数据库服务器配置
type ConfigDB struct {
	PprofIP  string
	ListenIP string
	GameTablePrefix	string
	LogTablePrefix string
}

// 全局配置信息
type GConfig struct {
	DB ConfigDB
	Login ConfigLogin
	Game ConfigGame
	Manager ConfigManager
	Proxy ConfigProxy
}
