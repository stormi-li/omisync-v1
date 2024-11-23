# Omisync 分布式锁框架
**作者**: stormi-li  
**Email**: 2785782829@qq.com  
## 简介

**Omisync** 是一个简化分布式锁使用的框架，使得分布式锁的操作像本地锁一样简单直观。它通过与 Redis 集成，提供高效的分布式锁机制，适用于分布式系统中多个进程或服务之间的同步和互斥。

## 功能

- **支持看门狗机制**：自动延长锁的过期时间，避免因进程宕机等原因导致锁被错误释放。
- **支持锁ID**：每个锁都可以有唯一的 ID，支持多个锁并发控制。
- **支持阻塞**：在锁被占用时，其他进程可以等待锁释放。
## 教程
### 安装
```shell
go get github.com/stormi-li/omisync-v1
```
### 使用
```go
package main

import (
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omisync-v1"
)

// 初始化 Omisync 客户端
var c = omisync.NewClient(&redis.Options{Addr: "localhost:6379"})

func main() {
	// 启动两个并发的任务，模拟多进程竞争锁
	go process("Process 1")
	go process("Process 2")
	select {} // 保持主线程运行
}

func process(name string) {
	// 创建一个锁实例
	lock := c.NewLock("my-distributed-lock")

	// 尝试获取锁
	lock.Lock()
	log.Println(name, "成功抢占锁")

	// 模拟业务逻辑执行
	time.Sleep(1 * time.Second)

	// 释放锁
	lock.Unlock()
	log.Println(name, "释放锁")
}
```