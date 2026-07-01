package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pcms/config"
	"pcms/utils"

	"github.com/gin-gonic/gin"
)

// Upload 文件上传
func Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请选择文件")
		return
	}
	defer file.Close()

	// 限制文件大小 (10MB)
	if header.Size > 10*1024*1024 {
		utils.BadRequest(c, "文件大小不能超过 10MB")
		return
	}

	// 生成唯一文件名
	ext := filepath.Ext(header.Filename)
	newFileName := fmt.Sprintf("%d_%s%s", time.Now().UnixMilli(),
		strings.ReplaceAll(strings.TrimSuffix(header.Filename, ext), " ", "_"),
		ext)

	// 按日期创建子目录
	dateDir := time.Now().Format("2006/01/02")
	uploadPath := filepath.Join(config.AppConfig.UploadDir, dateDir)
	os.MkdirAll(uploadPath, 0755)

	// 保存文件
	filePath := filepath.Join(uploadPath, newFileName)
	dst, err := os.Create(filePath)
	if err != nil {
		utils.InternalError(c, "文件保存失败")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		utils.InternalError(c, "文件写入失败")
		return
	}

	// 返回文件 URL（相对路径）
	relativePath := filepath.Join(dateDir, newFileName)
	relativePath = strings.ReplaceAll(relativePath, "\\", "/")

	utils.Success(c, gin.H{
		"url":       "/api/v1/files/" + relativePath,
		"file_name": header.Filename,
		"file_size": header.Size,
		"mime_type": http.DetectContentType([]byte{}),
	})
}

// ServeFile 静态文件服务
func ServeFile(c *gin.Context) {
	filePath := c.Param("filepath")
	fullPath := filepath.Join(config.AppConfig.UploadDir, filePath)

	// 安全检查：防止目录遍历
	if !strings.HasPrefix(filepath.Clean(fullPath), filepath.Clean(config.AppConfig.UploadDir)) {
		utils.Error(c, http.StatusForbidden, "禁止访问")
		return
	}

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		utils.NotFound(c, "文件不存在")
		return
	}

	c.File(fullPath)
}
