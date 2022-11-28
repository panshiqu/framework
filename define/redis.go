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
