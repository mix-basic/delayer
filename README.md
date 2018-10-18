# Delayer

高性能延迟队列中间件，采用 Golang 开发

## 如何使用

`delayer` 分为：

- 服务器端：负责定时扫描到期的任务，并放入队列，需在服务器上常驻执行。
- 客户端：在代码中使用，以类库的形式，可 `push`、`pop`、`remove` 任务。

## 服务器端

在 https://github.com/mixstart/delayer/releases 中下载对应平台的程序，解压后直接执行即可。

> 支持 windows、linux、mac 三种平台

启动：

```
[root@localhost bin]# ./delayerd
    ____       __
   / __ \___  / /___ ___  _____  _____
  / / / / _ \/ / __ `/ / / / _ \/ ___/
 / /_/ /  __/ / /_/ / /_/ /  __/ /
/_____/\___/_/\__,_/\__, /\___/_/
                   /____/
Service:		delayerd
Version:		1.0.1
[info] 2018/10/18 19:51:10 Service started successfully.
```

查看帮助：

```
[root@localhost bin]# ./delayerd -h
Usage: delayerd [options]

Options:
-c/--configuration FILENAME -- configuration file path (searches if not given)
-h/--help -- print this usage message and exit
-v/--version -- print version number and exit
```

## 客户端

我们提供了以下几种语言：

> 根据对应项目的说明使用

- PHP：https://github.com/mixstart/delayer-client-php
- Golang：开发中
- Java：待定
- Python：待定
