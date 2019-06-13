package utils

import (
	"gopkg.in/ini.v1"
	"log"
	"fmt"
)

// 配置数据
type Config struct {
	Delayer Delayer
	Redis   Redis
}

// delayer 节点数据
type Delayer struct {
	Pid               string
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
	// 默认文件
	if fileName == "" {
		fileName = "delayer.conf"
	}
	// 读取配置文件
	conf, err := ini.Load(fileName)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Configuration file read error: %s", fileName))
	}
	// 提取数据
	delayer := conf.Section("delayer")
	pid := delayer.Key("pid").String()
	timerInterval, _ := delayer.Key("timer_interval").Int64()
	accessLog := delayer.Key("access_log").String()
	errorLog := delayer.Key("error_log").String()
	redis := conf.Section("redis")
	host := redis.Key("host").String()
	port := redis.Key("port").String()
	database, _ := redis.Key("database").Int()
	password := redis.Key("password").String()
	maxIdle, _ := redis.Key("max_idle").Int()
	maxActive, _ := redis.Key("max_active").Int()
	idleTimeout, _ := redis.Key("idle_timeout").Int64()
	connMaxLifetime, _ := redis.Key("conn_max_lifetime").Int64()
	// 返回
	data := Config{
		Delayer: Delayer{
			Pid:           pid,
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
