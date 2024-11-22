package main

import (
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omisync-v1"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"
var c = omisync.NewClient(&redis.Options{Addr: redisAddr, Password: password})

func main() {
	go process("process1")
	go process("process2")
	select {}
}

func process(name string) {
	lock := c.NewLock("lock")
	lock.Lock()
	log.Println(name, "抢占锁")
	time.Sleep(1 * time.Second)
	lock.Unlock()
	log.Println(name, "释放锁")
}
