package cmd

import (
	"delayer/utils"
	"delayer/logic"
	"fmt"
)

var (
	config utils.Config
	logger utils.Logger
)

const (
	APP_VERSION = "1.0.1-dev"
)

func init() {
	config = utils.LoadConfig("delayer.conf")
	logger = utils.NewLogger(config)
}

func welcome() {
	fmt.Println("    ____       __                     ");
	fmt.Println("   / __ \\___  / /___ ___  _____  _____");
	fmt.Println("  / / / / _ \\/ / __ `/ / / / _ \\/ ___/");
	fmt.Println(" / /_/ /  __/ / /_/ / /_/ /  __/ /    ");
	fmt.Println("/_____/\\___/_/\\__,_/\\__, /\\___/_/     ");
	fmt.Println("                   /____/             ");
	fmt.Println("ServiceName:	delayerd");
	fmt.Println("Version:		7.2.9");
}

func Run() {
	// 欢迎
	welcome()
	// 启动定时器
	timer := logic.Timer{
		Config: config,
		Logger: logger,
	}
	timer.Init()
	timer.Run()
}
