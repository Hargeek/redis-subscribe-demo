package model

type Subscription struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    string `gorm:"index"` // 飞书用户ID
	Service   string `gorm:"index"` // 服务名称
	Branch    string `gorm:"index"` // 分支名称
	CreatedAt int64
	UpdatedAt int64
}

func (*Subscription) TableName() string {
	return "subscription"
}

type Notification struct {
	Service         string `json:"service"`
	Branch          string `json:"branch"`
	DefaultReceiver string `json:"default_receiver"` // 默认接收者
	Status          string `json:"status"`           // 构建状态
	Result          string `json:"result"`           // 构建结果
}
