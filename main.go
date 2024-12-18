package main

import (
	"context"
	"github.com/gin-gonic/gin"
	rediss "github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"redis-subscribe-demo/internal/config"
	"redis-subscribe-demo/internal/handler"
	"redis-subscribe-demo/internal/repository"
	"redis-subscribe-demo/internal/service"
	"redis-subscribe-demo/pkg/feishu"
	"redis-subscribe-demo/pkg/redis"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 使用配置初始化数据库
	db, err := gorm.Open(mysql.Open(cfg.MySQL.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 设置数据库连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)

	// 初始化Redis客户端
	redisClient, err := redis.NewRedisClient(&rediss.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// 初始化飞书客户端
	fsClient := feishu.NewClient()

	// 初始化仓储层
	repo := repository.NewSubscriptionRepo(db)

	// 初始化服务层
	subService := service.NewSubscriptionService(*repo, redisClient.GetClient(), fsClient)

	// 启动订阅监听
	go subService.StartSubscribe(context.Background())

	// 初始化处理器
	h := handler.NewHandler(subService)

	// 设置路由
	r := gin.Default()

	// API路由组
	v1 := r.Group("/api/v1")
	{
		v1.POST("/subscription", h.CreateSubscription)
		//v1.DELETE("/subscription", h.DeleteSubscription)
		v1.POST("/notification", h.HandleNotification)
		//v1.GET("/subscription", h.ListSubscriptions)
	}

	r.Run(cfg.Server.Addr)
}
