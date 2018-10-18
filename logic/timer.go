package logic

import (
	"delayer/utils"
	"github.com/gomodule/redigo/redis"
	"time"
	"strings"
	"fmt"
)

type Timer struct {
	Config    utils.Config
	Logger    utils.Logger
	Ticker    *time.Ticker
	Pool      *redis.Pool
	ErrHandle func(err error, funcName string, data string)
}

const (
	KEY_JOP_POOL       = "delayer:jop_pool"
	PREFIX_JOP_BUCKET  = "delayer:jop_bucket:"
	PREFIX_READY_QUEUE = "delayer:ready_queue:"
)

// 初始化
func (p *Timer) Init() {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", p.Config.Redis.Host+":"+p.Config.Redis.Port)
			if err != nil {
				return nil, err
			}
			if (p.Config.Redis.Password != "") {
				if _, err := c.Do("AUTH", p.Config.Redis.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if _, err := c.Do("SELECT", p.Config.Redis.Database); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		MaxIdle:         p.Config.Redis.MaxIdle,
		MaxActive:       p.Config.Redis.MaxActive,
		IdleTimeout:     time.Duration(p.Config.Redis.IdleTimeout) * time.Second,
		MaxConnLifetime: time.Duration(p.Config.Redis.ConnMaxLifetime) * time.Second,
	}
	p.Pool = pool
	errHandle := func(err error, funcName string, data string) {
		if (err != nil) {
			if (data != "") {
				data = ", [" + data + "]"
			}
			message := fmt.Sprintf("FAILURE: func %s, %s%s", funcName, err.Error(), data)
			p.Logger.Error(message)
		}
	}
	p.ErrHandle = errHandle
}

// 开始
func (p *Timer) Start() {
	ticker := time.NewTicker(time.Duration(p.Config.Delayerd.TimerInterval) * time.Millisecond)
	go func() {
		for range ticker.C {
			p.run()
		}
	}()
	p.Ticker = ticker
}

// 执行任务
func (p *Timer) run() {
	// 获取到期的任务
	jobs, err := p.getExpireJobs()
	if (err != nil) {
		p.ErrHandle(err, "getExpireJobs", "")
		return
	}
	// 并行获取Topic
	topics := make(map[string][]string)
	ch := make(chan []string)
	for _, jobID := range jobs {
		go p.getJopTopic(jobID, ch)
	}
	// Topic分组
	for i := 0; i < len(jobs); i++ {
		arr := <-ch
		if (arr[1] != "") {
			if _, ok := topics[arr[1]]; !ok {
				jobIDs := []string{arr[0]}
				topics[arr[1]] = jobIDs
			} else {
				topics[arr[1]] = append(topics[arr[1]], arr[0])
			}
		}
	}
	// 并行移动至Topic对应的ReadyQueue
	for topic, jobIDs := range topics {
		go p.moveJobToReadyQueue(jobIDs, topic)
	}
}

// 获取到期的任务
func (p *Timer) getExpireJobs() ([]string, error) {
	conn := p.Pool.Get()
	defer conn.Close()
	return redis.Strings(conn.Do("ZRANGEBYSCORE", KEY_JOP_POOL, "0", time.Now().Unix()))
}

// 获取任务的Topic
func (p *Timer) getJopTopic(jobID string, ch chan []string) {
	conn := p.Pool.Get()
	defer conn.Close()
	topic, err := redis.Strings(conn.Do("HMGET", PREFIX_JOP_BUCKET+jobID, "topic"))
	if (err != nil) {
		p.ErrHandle(err, "getJopTopic", jobID)
		ch <- []string{jobID, ""}
		return
	}
	arr := []string{jobID, topic[0]}
	ch <- arr
}

// 移动任务至ReadyQueue
func (p *Timer) moveJobToReadyQueue(jobIDs []string, topic string) {
	// 获取连接
	conn := p.Pool.Get()
	defer conn.Close()
	jobIDsStr := strings.Join(jobIDs, ",")
	// 开启事物
	if err := p.startTrans(conn); err != nil {
		p.ErrHandle(err, "startTrans", jobIDsStr)
		return
	}
	// 移除JopPool
	if err := p.delJopPool(conn, jobIDs, topic); err != nil {
		p.ErrHandle(err, "delJopPool", jobIDsStr)
		return
	}
	// 插入ReadyQueue
	if err := p.addReadyQueue(conn, jobIDs, topic); err != nil {
		p.ErrHandle(err, "addReadyQueue", jobIDsStr)
		return
	}
	// 提交事物
	if err := p.commit(conn); err != nil {
		p.ErrHandle(err, "commit", jobIDsStr)
		return
	}
	// 打印日志
	message := fmt.Sprintf("Job is ready, Topic: %s, IDs: [%s]", topic, jobIDsStr)
	p.Logger.Info(message)
}

// 开启事务
func (p *Timer) startTrans(conn redis.Conn) error {
	return conn.Send("MULTI")
}

// 提交事务
func (p *Timer) commit(conn redis.Conn) error {
	return conn.Send("EXEC")
}

// 移除JopPool
func (p *Timer) delJopPool(conn redis.Conn, jobIDs []string, topic string) error {
	args := make([]interface{}, len(jobIDs)+1)
	args[0] = KEY_JOP_POOL
	for k, v := range jobIDs {
		args[k+1] = v
	}
	return conn.Send("ZREM", args...)
}

// 插入ReadyQueue
func (p *Timer) addReadyQueue(conn redis.Conn, jobIDs []string, topic string) error {
	args := make([]interface{}, len(jobIDs)+1)
	args[0] = PREFIX_READY_QUEUE + topic
	for k, v := range jobIDs {
		args[k+1] = v
	}
	return conn.Send("LPUSH", args...)
}

// 执行
func (p *Timer) Stop() {
	p.Ticker.Stop()
}
