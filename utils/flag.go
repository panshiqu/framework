package utils

import (
	"../define"
	"flag"
)

//解析通用参数
func getCommonArgs(args *define.CommonArgs)  {
	flag.StringVar(& args.ConfigPath, "config", "./config/login.json", "config log path")
}

//解析登录参数
func GetLoginArgs() define.LoginArgs  {
	var args define.LoginArgs
	getCommonArgs(& args.CommonArgs)

	flag.Parse()
	return args
}

//解析数据库参数
func GetDBArgs() define.DBArgs {
	var args define.DBArgs
	getCommonArgs(& args.CommonArgs)

	flag.Parse()
	return  args
}