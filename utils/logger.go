package utils

import (
	"log"
	"os"
	"io"
	"fmt"
)

// 日志类
type Logger struct {
	AccessLog string
	ErrorLog  string
}

// 打开文件
func (p *Logger) openFile(fileName string) *os.File {
	logFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Open file Failed: %s", fileName))
	}
	return logFile
}

// 信息日志
func (p *Logger) Info(message string) {
	fileName := p.AccessLog
	var out io.Writer
	if (fileName == "") {
		out = os.Stdout
	} else {
		logFile := p.openFile(fileName)
		defer logFile.Close()
		out = io.MultiWriter(os.Stdout, logFile)
	}
	logLogger := log.New(out, "[info] ", log.LstdFlags)
	logLogger.Println(message)
}

// 错误日志
func (p *Logger) Error(message string, exit bool) {
	fileName := p.ErrorLog
	var out io.Writer
	if (fileName == "") {
		out = os.Stdout
	} else {
		logFile := p.openFile(fileName)
		defer logFile.Close()
		out = io.MultiWriter(logFile, os.Stdout)
	}
	logLogger := log.New(out, "[error] ", log.LstdFlags)
	if exit {
		logLogger.Fatalln(message)
	}
	logLogger.Println(message)
}

// 创建实例
func NewLogger(config Config) Logger {
	logger := Logger{
		AccessLog: config.Delayer.AccessLog,
		ErrorLog:  config.Delayer.ErrorLog,
	}
	return logger
}
