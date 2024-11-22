package omisync

import (
	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omipc-v1"
)

func NewClient(opts *redis.Options) *Client {
	return &Client{redisClient: redis.NewClient(opts), omipcClient: omipc.NewClient(opts)}
}
