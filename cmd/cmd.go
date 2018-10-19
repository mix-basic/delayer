package cmd

import (
	"delayer/utils"
	"delayer/logic"
	"fmt"
	"os"
	"flag"
	"os/signal"
	"syscall"
	"io/ioutil"
	"log"
)

const (
	APP_VERSION = "1.0.1"
)

func Run() {
	// 守护进程
	utils.Daemon()
	// 命令行参数处理
	configuration := handleFlag()
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
	// pid处理
	handlePid(config);
	// 输出启动日志
	message := fmt.Sprintf("Service started successfully, PID: %d", os.Getpid())
	logger.Info(message)
	// 启动定时器
	timer := logic.Timer{
		Config: config,
		Logger: logger,
	}
	timer.Init()
	timer.Start()
	// 信号处理
	handleSignal(timer, exit)
	// 退出
	<-exit
	// 输出停止日志
	message = fmt.Sprintf("Service stopped successfully, PID: %d", os.Getpid())
	logger.Info(message)
}

func welcome() {
	fmt.Println("    ____       __                     ");
	fmt.Println("   / __ \\___  / /___ ___  _____  _____");
	fmt.Println("  / / / / _ \\/ / __ `/ / / / _ \\/ ___/");
	fmt.Println(" / /_/ /  __/ / /_/ / /_/ /  __/ /    ");
	fmt.Println("/_____/\\___/_/\\__,_/\\__, /\\___/_/     ");
	fmt.Println("                   /____/             ");
	fmt.Println("Service:		delayer");
	fmt.Println("Version:		" + APP_VERSION);
}

func handlePid(config utils.Config) {
	// 读取
	pidStr, err := ioutil.ReadFile(config.Delayer.Pid)
	if err != nil {
		writePidFile(config.Delayer.Pid)
		return
	}
	// 重复启动处理
	pid, err := utils.ByteToInt(pidStr)
	if (err != nil) {
		writePidFile(config.Delayer.Pid)
		return
	}
	pro, err := os.FindProcess(pid)
	if err != nil {
		writePidFile(config.Delayer.Pid)
		return
	}
	err = pro.Signal(os.Signal(syscall.Signal(0)))
	if err != nil {
		// os: process already finished
		// not supported by windows
		writePidFile(config.Delayer.Pid)
		return
	}
	log.Fatalln(fmt.Sprintf("ERROR: Service is being executed, PID: %d", pid))
}

func writePidFile(pidFile string) {
	err := ioutil.WriteFile(pidFile, utils.IntToByte(os.Getpid()), 0644)
	if err != nil {
		log.Fatalln(fmt.Sprintf("PID file cannot be written: %s", pidFile))
	}
}

// 信号处理
func handleSignal(timer logic.Timer, exit chan bool) {
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
}

func handleFlag() string {
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
	fmt.Println()
	os.Exit(0)
}

func printVersion() {
	fmt.Println(APP_VERSION);
	fmt.Println()
	os.Exit(0)
}
