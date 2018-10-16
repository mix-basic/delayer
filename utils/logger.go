package utils

import (
	"log"
	"os"
	"io"
)

// 日志类
type Logger struct {
	AccessLog string
	ErrorLog  string
}

// 打开文件
func (p *Logger) openFile(fileName string) *os.File {
	logFileOpen, err := os.OpenFile(fileName, os.O_APPEND, 0644)
	if err == nil {
		return logFileOpen
	}
	logFileCreate, err := os.Create(fileName)
	if err != nil {
		log.Fatalln("Failed to create log file: " + fileName)
	}
	return logFileCreate
}

// 信息日志
func (p *Logger) Info(message string) {
	fileName := p.AccessLog
	var out io.Writer
	if (fileName == "") {
		out = os.Stdout
	} else {
		logFile := p.openFile(fileName)
		out = io.MultiWriter(logFile, os.Stdout)
	}
	logLogger := log.New(out, "[info] ", log.LstdFlags)
	logLogger.Println(message)
}

// 错误日志
// 会退出程序
func (p *Logger) Error(message string) {
	fileName := p.ErrorLog
	var out io.Writer
	if (fileName == "") {
		out = os.Stdout
	} else {
		logFile := p.openFile(fileName)
		out = io.MultiWriter(logFile, os.Stdout)
	}
	logLogger := log.New(out, "[error] ", log.LstdFlags)
	logLogger.Println(message)
}

// 创建实例
func NewLogger(config Config) Logger {
	logger := Logger{
		AccessLog: config.Delayerd.AccessLog,
		ErrorLog:  config.Delayerd.ErrorLog,
	}
	return logger
}
