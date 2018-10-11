package logger

import (
	"delayer/pkg/config"
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
	logFile := p.openFile(fileName)
	logLogger := log.New(io.MultiWriter(logFile, os.Stdout), "[info] ", log.LstdFlags)
	logLogger.Println(message)
}

// 错误日志
// 会退出程序
func (p *Logger) Error(message string) {
	fileName := p.ErrorLog
	logFile := p.openFile(fileName)
	logLogger := log.New(io.MultiWriter(logFile, os.Stderr), "[error] ", log.LstdFlags)
	logLogger.Fatalln(message)
}

// 创建实例
func New(settings config.Settings) Logger {
	logger := Logger{
		AccessLog: settings.Delayerd.AccessLog,
		ErrorLog:  settings.Delayerd.ErrorLog,
	}
	return logger
}
