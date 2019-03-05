package define

//通用参数
type CommonArgs struct {
	ConfigPath string
}

//登录参数
type LoginArgs struct {
	CommonArgs
}

//数据库参数
type DBArgs struct {
	CommonArgs
} 