package routes

import (
	"pcms/handlers"
	"pcms/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authHandler *handlers.AuthHandler,
	categoryHandler *handlers.CategoryHandler,
	documentHandler *handlers.DocumentHandler,
	tagHandler *handlers.TagHandler,
	attachmentHandler *handlers.AttachmentHandler,
) *gin.Engine {
	r := gin.Default()

	// CORS
	r.Use(middleware.CORS())

	// 静态文件服务
	r.GET("/api/v1/files/*filepath", handlers.ServeFile)

	api := r.Group("/api/v1")
	{
		// 认证（不需要JWT）
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// 需要JWT认证的路由
		protected := api.Group("")
		protected.Use(middleware.AuthRequired())
		{
			// 用户
			protected.GET("/auth/me", authHandler.Me)

			// 分类
			protected.GET("/categories", categoryHandler.GetTree)
			protected.POST("/categories", categoryHandler.Create)
			protected.PUT("/categories/:id", categoryHandler.Update)
			protected.DELETE("/categories/:id", categoryHandler.Delete)
			protected.PUT("/categories/:id/move", categoryHandler.Move)

			// 文档
			protected.GET("/documents", documentHandler.List)
			protected.GET("/documents/search", documentHandler.Search)
			protected.POST("/documents", documentHandler.Create)
			protected.GET("/documents/:id", documentHandler.Get)
			protected.PUT("/documents/:id", documentHandler.Update)
			protected.DELETE("/documents/:id", documentHandler.Delete)

			// 文档版本
			protected.GET("/documents/:id/versions", documentHandler.GetVersions)
			protected.GET("/documents/:id/versions/:vid", documentHandler.GetVersion)

			// 标签
			protected.GET("/tags", tagHandler.List)
			protected.POST("/tags", tagHandler.Create)
			protected.DELETE("/tags/:id", tagHandler.Delete)

			// 附件
			protected.GET("/attachments", attachmentHandler.List)
			protected.POST("/attachments", attachmentHandler.Create)
			protected.DELETE("/attachments/:id", attachmentHandler.Delete)
			protected.PUT("/attachments/:id/bind", attachmentHandler.BindDocument)

			// 文件上传
			protected.POST("/files/upload", handlers.Upload)
		}
	}

	return r
}
