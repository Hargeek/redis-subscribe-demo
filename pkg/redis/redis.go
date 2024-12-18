package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(opts *redis.Options) (*RedisClient, error) {
	client := redis.NewClient(opts)
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis连接失败: %w", err)
	}

	return &RedisClient{client: client}, nil
}

// GetClient 获取 Redis 客户端
func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

// Publish 发布消息
func (r *RedisClient) Publish(ctx context.Context, channel string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("消息序列化失败: %w", err)
	}
	return r.client.Publish(ctx, channel, data).Err()
}

// Subscribe 订阅消息
func (r *RedisClient) Subscribe(ctx context.Context, channel string, handler func([]byte)) error {
	sub := r.client.Subscribe(ctx, channel)
	defer sub.Close()

	ch := sub.Channel()
	for msg := range ch {
		handler([]byte(msg.Payload))
	}
	return nil
}
