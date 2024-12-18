package repository

import (
	"context"
	"gorm.io/gorm"
	"redis-subscribe-demo/internal/model"
)

type SubscriptionRepo struct {
	DB *gorm.DB
}

func NewSubscriptionRepo(db *gorm.DB) *SubscriptionRepo {
	return &SubscriptionRepo{DB: db}
}

func (r *SubscriptionRepo) Create(ctx context.Context, sub *model.Subscription) error {
	return r.DB.WithContext(ctx).Create(sub).Error
}

func (r *SubscriptionRepo) GetSubscriptions(ctx context.Context, serviceName, branchName string) ([]model.Subscription, error) {
	var subs []model.Subscription
	err := r.DB.WithContext(ctx).
		Where("service = ? AND branch = ?", serviceName, branchName).
		Find(&subs).Error
	return subs, err
}
