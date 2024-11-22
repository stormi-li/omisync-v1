package omisync

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stormi-li/omipc-v1"
)

type Client struct {
	redisClient *redis.Client
	omipcClient *omipc.Client
}

func (c *Client) NewLock(lockName string) *Lock {
	return &Lock{
		uuid:        uuid.NewString(),
		lockName:    lockName,
		stop:        make(chan struct{}, 1),
		omipcClient: c.omipcClient,
		redisClient: c.redisClient,
		ctx:         context.Background(),
	}
}

func (c *Client) Close() {
	c.redisClient.Close()
	c.omipcClient.Close()
}
