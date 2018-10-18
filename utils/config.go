package utils

import (
	"gopkg.in/ini.v1"
	"log"
	"fmt"
)

// 配置数据
type Config struct {
	Delayerd Delayerd
	Redis    Redis
}

// delayerd 节点数据
type Delayerd struct {
	TimerInterval     int64
	BucketMaxLifetime int64
	AccessLog         string
	ErrorLog          string
}

// redis 节点数据
type Redis struct {
	Host            string
	Port            string
	Database        int
	Password        string
	MaxIdle         int
	MaxActive       int
	IdleTimeout     int64
	ConnMaxLifetime int64
}

// 载入配置
func LoadConfig(fileName string) Config {
	// 读取配置文件
	conf, err := ini.Load(fileName)
	if err != nil {
		message := fmt.Sprintf("Configuration file read error: %s", fileName)
		log.Fatalln(message)
	}
	// 提取数据
	delayerd := conf.Section("delayerd")
	timerInterval, _ := delayerd.Key("timer_interval").Int64()
	accessLog := delayerd.Key("access_log").String()
	errorLog := delayerd.Key("error_log").String()
	redis := conf.Section("redis")
	host := redis.Key("host").String()
	port := redis.Key("port").String()
	database, _ := delayerd.Key("database").Int()
	password := redis.Key("password").String()
	maxIdle, _ := delayerd.Key("max_idle").Int()
	maxActive, _ := delayerd.Key("max_active").Int()
	idleTimeout, _ := delayerd.Key("idle_timeout").Int64()
	connMaxLifetime, _ := delayerd.Key("conn_max_lifetime").Int64()
	// 返回
	data := Config{
		Delayerd: Delayerd{
			TimerInterval: timerInterval,
			AccessLog:     accessLog,
			ErrorLog:      errorLog,
		},
		Redis: Redis{
			Host:            host,
			Port:            port,
			Database:        database,
			Password:        password,
			MaxIdle:         maxIdle,
			MaxActive:       maxActive,
			IdleTimeout:     idleTimeout,
			ConnMaxLifetime: connMaxLifetime,
		},
	}
	return data
}
