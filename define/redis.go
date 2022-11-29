package define

const (
	// RedisOnline 在线
	RedisOnline = iota
)

// Database 0

// UserID 用户编号
// GameID 游戏服务编号
// 标识用户在线哪个游戏服务
// KEY: Online_GameID_UserID, Example: Online_2_1
// Value: JSON, Example: {"UserID":1,"GameID":2,"GameType":2,"GameLevel":1}
// 当前实现数据结构和逻辑简单粗暴，若未来数据量多到KEYS命令存在效率问题，则可拆分成如下结构来实现
// Hash Online_GameID-UserID / Online_UserID-JSON
