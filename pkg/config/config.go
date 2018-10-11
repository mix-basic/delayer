package config

import (
	"gopkg.in/ini.v1"
	"log"
)

// 配置数据
type Settings struct {
	Delayerd Delayerd
	Redis    Redis
}

// delayerd 节点数据
type Delayerd struct {
	Interval  int
	AccessLog string
	ErrorLog  string
}

// redis 节点数据
type Redis struct {
	Host      string
	Port      string
	Database  int
	Password  string
	Timeout   int
	MaxIdle   int
	MaxActive int
}

// 载入配置
func Load(fileName string) Settings {
	// 读取配置文件
	conf, err := ini.Load(fileName)
	if err != nil {
		log.Fatalln("Configuration file read error: " + fileName)
	}
	// 提取数据
	delayerd := conf.Section("delayerd")
	interval, _ := delayerd.Key("interval").Int()
	accessLog := delayerd.Key("access_log").String()
	errorLog := delayerd.Key("error_log").String()
	redis := conf.Section("redis")
	host := redis.Key("host").String()
	port := redis.Key("port").String()
	database, _ := delayerd.Key("database").Int()
	password := redis.Key("password").String()
	timeout, _ := delayerd.Key("timeout").Int()
	maxIdle, _ := delayerd.Key("max_idle").Int()
	maxActive, _ := delayerd.Key("max_active").Int()
	// 返回
	data := Settings{
		Delayerd: Delayerd{
			Interval:  interval,
			AccessLog: accessLog,
			ErrorLog:  errorLog,
		},
		Redis: Redis{
			Host:      host,
			Port:      port,
			Database:  database,
			Password:  password,
			Timeout:   timeout,
			MaxIdle:   maxIdle,
			MaxActive: maxActive,
		},
	}
	return data
}
