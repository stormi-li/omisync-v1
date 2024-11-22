package omisync

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omipc-v1"
)

type Lock struct {
	uuid        string
	lockName    string
	stop        chan struct{}
	redisClient *redis.Client
	omipcClient *omipc.Client
	ctx         context.Context
}

func (l *Lock) Lock() {
	for {
		var ok bool
		//尝试占有锁-----------------------------------------redis代码
		ok, _ = l.redisClient.SetNX(l.ctx, l.lockName, l.uuid, lock_expire_time*time.Second).Result()

		if ok {
			//看门口协程
			go func() {
				ticker := time.NewTicker(watchdog_interval)
				defer ticker.Stop()
				for {
					select {
					case <-ticker.C:
						//如果占有锁则更新过期时间
						l.updateExpiryIfValueMatches()
					case <-l.stop:
						return
					}
				}
			}()
			break
		} else {
			//阻塞三秒，阻塞时可以被唤醒
			l.omipcClient.Listen(l.lockName, block_wait_time)
		}
	}
}

func (l *Lock) Unlock() {
	l.stop <- struct{}{}
	l.deleteIfValueMatches()
	l.omipcClient.Notify(l.lockName, "unlock")

}

func (l *Lock) IsValid() bool {
	res, _ := l.updateExpiryIfValueMatches()
	return res
}

// updateExpiryIfValueMatches 更新过期时间，如果键的值匹配预期值
func (l *Lock) updateExpiryIfValueMatches() (bool, error) {
	script := `
        local key = KEYS[1]
        local expected_value = ARGV[1]
        local new_ttl = ARGV[2]

        local current_value = redis.call('GET', key)

        if current_value == expected_value then
            redis.call('EXPIRE', key, new_ttl)
            return 1
        else
            return 0
        end
    `
	//执行lua脚本-----------------------------------------redis代码
	result, err := l.redisClient.Eval(l.ctx, script, []string{l.lockName}, l.uuid, lock_expire_time).Result()
	if err != nil {
		return false, err
	}

	return result.(int64) == 1, nil
}

// deleteIfValueMatches 判断键的值是否匹配目标值，若匹配则删除该键
func (l *Lock) deleteIfValueMatches() (bool, error) {
	luaScript := `
		local currentValue = redis.call("GET", KEYS[1])
		if currentValue == ARGV[1] then
			redis.call("DEL", KEYS[1])
			return 1  -- 1表示成功删除
		else
			return 0  -- 0表示值不匹配
		end
	`
	//执行lua脚本-----------------------------------------redis代码
	result, err := l.redisClient.Eval(l.ctx, luaScript, []string{l.lockName}, l.uuid).Result()
	if err != nil {
		return false, err
	}
	return result.(int64) == 1, nil
}
