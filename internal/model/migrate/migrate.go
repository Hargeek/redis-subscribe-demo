package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"redis-subscribe-demo/internal/config"
	"redis-subscribe-demo/internal/model"
)

func main() {
	cfg := config.LoadConfig()
	db, err := gorm.Open(mysql.Open(cfg.MySQL.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)
	if err = db.Debug().AutoMigrate(
		&model.Subscription{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrated successfully")
}
