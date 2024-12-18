package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"redis-subscribe-demo/internal/model"
	"redis-subscribe-demo/internal/repository"
	"redis-subscribe-demo/pkg/feishu"
)

type SubscriptionService struct {
	repo     repository.SubscriptionRepo
	rdb      *redis.Client
	fsClient *feishu.Client
}

// NewSubscriptionService 创建订阅服务实例
func NewSubscriptionService(repo repository.SubscriptionRepo, redisClient *redis.Client, fsClient *feishu.Client) *SubscriptionService {
	return &SubscriptionService{
		repo:     repo,
		rdb:      redisClient,
		fsClient: fsClient,
	}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, userID, service, branch string) error {
	sub := &model.Subscription{
		UserID:  userID,
		Service: service,
		Branch:  branch,
	}

	if err := s.repo.Create(ctx, sub); err != nil {
		return fmt.Errorf("保存订阅到数据库失败: %w", err)
	}

	// 3. 更新 Redis 缓存
	key := fmt.Sprintf("subscription:%s:%s", service, branch)
	if err := s.rdb.HSet(ctx, key, userID, "1").Err(); err != nil {
		log.Printf("更新Redis缓存失败: %v", err)
	}

	return nil
}

func (s *SubscriptionService) HandleNotification(ctx context.Context, notification *model.Notification) error {
	// 构建通知消息
	channel := fmt.Sprintf("pipeline:%s:%s", notification.Service, notification.Branch)

	// 发布通知消息
	if err := s.rdb.Publish(ctx, channel, notification).Err(); err != nil {
		return fmt.Errorf("发布通知失败: %w", err)
	}

	return nil
}

// StartSubscribe 启动订阅监听
func (s *SubscriptionService) StartSubscribe(ctx context.Context) {
	// 可以根据需要订阅不同的通道
	channels := []string{"pipeline:*"}

	for _, channel := range channels {
		go func(ch string) {
			pubSub := s.rdb.Subscribe(ctx, ch)
			defer pubSub.Close()

			for msg := range pubSub.Channel() {
				var notification model.Notification
				if err := json.Unmarshal([]byte(msg.Payload), &notification); err != nil {
					log.Printf("解析通知消息失败: %v", err)
					continue
				}

				s.fsClient.SendMessage(s.fsClient.BuildNotificationCard(&notification), []string{notification.DefaultReceiver})
			}
		}(channel)
	}
}
