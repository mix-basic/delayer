package main

import (
	"delayer/pkg/config"
	"delayer/pkg/logger"
)

var (
	Settings config.Settings
	Logger   logger.Logger
)

func init() {
	Settings = config.Load("delayer.conf")
	Logger = logger.New(Settings)
}

func main() {

}
