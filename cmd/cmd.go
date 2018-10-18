package cmd

import (
	"delayer/utils"
	"delayer/logic"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"flag"
)

const (
	APP_VERSION = "1.0.1"
)

func Run() {
	// 命令行参数处理
	flagHandle();
	// 变量定义
	exit := make(chan bool)
	// 欢迎
	welcome()
	// 实例化公共组件
	config := utils.LoadConfig("delayer.conf")
	logger := utils.NewLogger(config)
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

func flagHandle() {
	h := flag.Bool("h", false, "")
	help := flag.Bool("help", false, "")
	v := flag.Bool("v", false, "")
	version := flag.Bool("version", false, "")
	flag.Parse()
	if *h || *help {
		printHelp()
	}
	if *v || *version {
		printVersion()
	}
}

func printHelp() {
	fmt.Println("Usage: delayerd [options]");
	fmt.Println()
	fmt.Println("Options:");
	fmt.Println("-c/--configuration FILENAME -- configuration file path (searches if not given)");
	fmt.Println("-h/--help -- print this usage message and exit");
	fmt.Println("-v/--version -- print version number and exit");
	os.Exit(0)
}

func printVersion() {
	fmt.Println(APP_VERSION);
	os.Exit(0)
}
