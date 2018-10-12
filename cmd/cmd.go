package cmd

import (
	"delayer/utils"
	"delayer/logic"
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

func Run() {
	// 连接
	timer := logic.Timer{
		Config: config,
		Logger: logger,
	}
	timer.Init()
	timer.Run()
}
