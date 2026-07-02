package main

import (
	"fmt"
	"log"

	"pcms/config"
	"pcms/handlers"
	"pcms/models"
	"pcms/routes"
	"pcms/services"
	"pcms/store"
	"pcms/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg := config.Load()
	utils.InitJWT(cfg.JWTSecret)

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Document{},
		&models.Tag{},
		&models.Attachment{},
		&models.DocumentVersion{},
	); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	st := store.NewGormStore(db)

	authHandler := handlers.NewAuthHandler(services.NewAuthService(st))
	categoryHandler := handlers.NewCategoryHandler(services.NewCategoryService(st))
	documentHandler := handlers.NewDocumentHandler(services.NewDocumentService(st))
	tagHandler := handlers.NewTagHandler(services.NewTagService(st))
	attachmentHandler := handlers.NewAttachmentHandler(services.NewAttachmentService(st))

	r := routes.SetupRouter(authHandler, categoryHandler, documentHandler, tagHandler, attachmentHandler)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("PCMS 服务启动在 http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
