# Delayer

基于 Redis 的延迟消息队列中间件，采用 Golang 开发，支持 PHP、Golang 等多种语言客户端。

参考 [有赞延迟队列设计](http://tech.youzan.com/queuing_delay) 中的部分设计，优化后实现。

## 应用场景

- 订单超过30分钟未支付，自动关闭订单。
- 订单完成后, 如果用户一直未评价, 5天后自动好评。
- 会员到期前3天，短信通知续费。
- 其他针对某个任务，延迟执行功能的需求。

## 实现原理

- 客户端：push 任务时，任务数据存入 hash 中，jobID 存入 zset 中，pop 时从指定的 list 中取准备好的数据。
- 服务器端：定时使用连接池并行将 zset 中到期的 jobID 放入对应的 list 中，供客户端 pop 取出。

## 核心特征

- 使用 Golang 开发，高性能。
- 高可用：服务器端操作是原子的，并且做了优雅停止，不会丢失数据，在redis断线时会自动重连。
- 可通过配置文件控制执行性能参数。
- 提供多种语言的 SDK，使用简单快捷。

## 如何使用

`delayer` 分为：

- 服务器端：负责定时扫描到期的任务，并放入队列，需在服务器上常驻执行。
- 客户端：在代码中使用，以类库的形式，提供 `push`、`pop`、`bPop`、`remove` 方法操作任务。

## 服务器端

在 https://github.com/mix-start/delayer/releases 中下载对应平台的程序。

> 支持 windows、linux、mac 三种平台

然后修改配置文件 `delayer.conf`：

```
[delayer]
pid = /var/run/delayer.pid      ; 需单例执行时配置, 多实例执行时留空, Win不支持单例
timer_interval = 1000           ; 计算间隔时间, 单位毫秒
access_log = logs/access.log    ; 存取日志
error_log = logs/error.log      ; 错误日志

[redis]
host = 127.0.0.1                ; 连接地址
port = 6379                     ; 连接端口
database = 0                    ; 数据库编号
password =                      ; 密码, 无需密码留空
max_idle = 2                    ; 最大空闲连接数
max_active = 20                 ; 最大激活连接数
idle_timeout = 3600             ; 空闲连接超时时间, 单位秒
conn_max_lifetime = 3600        ; 连接最大生存时间, 单位秒
```

查看帮助：

```
[root@localhost bin]# ./delayer -h
Usage: delayer [options]

Options:
-d/--daemon run in the background
-c/--configuration FILENAME -- configuration file path (searches if not given)
-h/--help -- print this usage message and exit
-v/--version -- print version number and exit
```

启动：

```
[root@localhost bin]# ./delayer
    ____       __
   / __ \___  / /___ ___  _____  _____
  / / / / _ \/ / __ `/ / / / _ \/ ___/
 / /_/ /  __/ / /_/ / /_/ /  __/ /
/_____/\___/_/\__,_/\__, /\___/_/
                   /____/
Service:		delayer
Version:		1.0.1
[info] 2018/10/21 11:24:24 Service started successfully, PID: 31023
```

## 客户端

我们提供了以下几种语言：

> 根据对应项目的说明使用

- PHP：https://github.com/mix-basic/delayer-client-php
- Golang：https://github.com/mix-basic/delayer-client-golang
- Java：待定
- Python：待定

## License

Apache License Version 2.0, http://www.apache.org/licenses/
