package utils

import (
	"../define"
	"flag"
)

func GetLoginArgs() define.LoginArgs  {
	var args define.LoginArgs
	flag.StringVar(&args.ConfigPath, "config", "./config/login.json", "config log path")

	flag.Parse()
	return args
}