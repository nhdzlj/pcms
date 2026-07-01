package main

import (
	"fmt"
	"log"

	"pcms/config"
	"pcms/handlers"
	"pcms/models"
	"pcms/routes"
	"pcms/services"
	"pcms/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化 JWT
	utils.InitJWT(cfg.JWTSecret)

	// 连接数据库
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 自动迁移
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

	// 初始化 Service
	authService := services.NewAuthService(db)
	categoryService := services.NewCategoryService(db)
	documentService := services.NewDocumentService(db)
	tagService := services.NewTagService(db)

	// 初始化 Handler
	authHandler := handlers.NewAuthHandler(authService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	documentHandler := handlers.NewDocumentHandler(documentService)
	tagHandler := handlers.NewTagHandler(tagService)

	// 设置路由
	r := routes.SetupRouter(authHandler, categoryHandler, documentHandler, tagHandler)

	// 启动服务
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("PCMS 服务启动在 http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
