package feishu

import (
	"log"
	"redis-subscribe-demo/internal/model"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

type MessageCard struct {
	Title   string
	Content string
	Status  string
}

// SendMessage 模拟发送消息，实际只打印日志
func (c *Client) SendMessage(notification *MessageCard, receivers []string) {
	log.Printf("=== Pipeline Notification ===\n"+
		"Service: %v\n"+
		"Branch: %v\n"+
		"Status: %v\n"+
		"Receivers: %v\n"+
		"========================",
		notification.Title,
		notification.Content,
		notification.Status,
		receivers)
}

// BuildNotificationCard 构建飞书通知卡片
func (c *Client) BuildNotificationCard(notification *model.Notification) *MessageCard {
	return &MessageCard{
		Title:   notification.Service + " 流水线通知",
		Content: notification.Branch + " 分支构建" + notification.Status,
		Status:  notification.Status,
	}
}
