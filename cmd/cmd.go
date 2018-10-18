package cmd

import (
	"delayer/utils"
	"delayer/logic"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	config utils.Config
	logger utils.Logger
)

const (
	APP_VERSION = "1.0.1"
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
	fmt.Println("Service:		delayerd");
	fmt.Println("Version:		" + APP_VERSION);
}

func Run() {
	exit := make(chan bool)
	// 欢迎
	welcome()
	// 启动定时器
	timer := logic.Timer{
		Config: config,
		Logger: logger,
	}
	timer.Init()
	timer.Start()
	// 信号处理
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-ch
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			timer.Stop()
			exit <- true
		}
	}()
	// 退出
	<-exit
}
