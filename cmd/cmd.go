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
	configuration := flagHandle()
	// 变量定义
	exit := make(chan bool)
	// 欢迎
	welcome()
	// 实例化公共组件
	if configuration == "" {
		configuration = "delayer.conf"
	}
	config := utils.LoadConfig(configuration)
	logger := utils.NewLogger(config)
	// 输出启动日志
	logger.Info("Service started successfully.")
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
	// 输出停止日志
	logger.Info("Service stopped successfully.")
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

func flagHandle() string {
	// 参数解析
	flagH := flag.Bool("h", false, "")
	flagHelp := flag.Bool("help", false, "")
	flagV := flag.Bool("v", false, "")
	flagVersion := flag.Bool("version", false, "")
	flagC := flag.String("c", "", "")
	flagConfiguration := flag.String("configuration", "", "")
	flag.Parse()
	// 参数取值
	help := *flagH || *flagHelp
	version := *flagV || *flagVersion
	configuration := ""
	if (*flagC == "") {
		configuration = *flagConfiguration
	} else {
		configuration = *flagC
	}
	// 打印型命令处理
	if help {
		printHelp()
	}
	if version {
		printVersion()
	}
	// 返回参数值
	return configuration
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
