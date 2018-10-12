package logic

import (
	"delayer/utils"
	"github.com/gomodule/redigo/redis"
	"time"
	"strings"
)

type Timer struct {
	Config utils.Config
	Logger utils.Logger
	Pool   *redis.Pool
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
}

// 执行
func (p *Timer) Run() {
	ticker := time.NewTicker(time.Duration(p.Config.Delayerd.TimerInterval) * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			// 执行任务
			jobs := p.getExpireJobs()
			// 并行获取Topic
			topics := make(map[string][]string)
			ch := make(chan []string)
			for _, jobID := range jobs {
				go p.getJopTopic(jobID, ch)
			}
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
	}
}

// 获取到期的任务
func (p *Timer) getExpireJobs() []string {
	conn := p.Pool.Get()
	defer conn.Close()
	jobs, err := redis.Strings(conn.Do("ZRANGEBYSCORE", KEY_JOP_POOL, "0", time.Now().Unix()))
	if (err != nil) {
		p.Logger.Error(err.Error())
	}
	return jobs
}

// 放入Ready队列
func (p *Timer) getJopTopic(jobID string, ch chan []string) {
	conn := p.Pool.Get()
	defer conn.Close()
	topic, err := redis.Strings(conn.Do("HMGET", PREFIX_JOP_BUCKET+jobID, "topic"))
	if (err != nil) {
		p.Logger.Error(err.Error())
	}
	arr := []string{jobID, topic[0]}
	ch <- arr
}

// 移动任务至ReadyQueue
func (p *Timer) moveJobToReadyQueue(jobIDs []string, topic string) {
	conn := p.Pool.Get()
	defer conn.Close()
	// 移除JopPool
	zremArgs := make([]interface{}, len(jobIDs)+1)
	zremArgs[0] = KEY_JOP_POOL
	for k, v := range jobIDs {
		zremArgs[k+1] = v
	}
	succ, err := redis.Bool(conn.Do("ZREM", zremArgs...))
	if (err != nil) {
		p.Logger.Error(err.Error())
	}
	if (!succ) {
		p.Logger.Info("Move failure, jobIDs: " + strings.Join(jobIDs, ","))
		return
	}
	// 插入ReadyQueue
	lpushArgs := make([]interface{}, len(jobIDs)+1)
	lpushArgs[0] = PREFIX_READY_QUEUE+topic
	for k, v := range jobIDs {
		lpushArgs[k+1] = v
	}
	succ, err := redis.Bool(conn.Do("LPUSH", zremArgs...))
	if (err != nil) {
		p.Logger.Error(err.Error())
	}
	if (!succ) {
		p.Logger.Info("Move failure, jobIDs: " + strings.Join(jobIDs, ","))
	}
}
