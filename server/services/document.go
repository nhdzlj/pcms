package services

import (
	"errors"
	"fmt"

	"pcms/models"
	"pcms/utils"

	"gorm.io/gorm"
)

type DocumentService struct {
	DB *gorm.DB
}

func NewDocumentService(db *gorm.DB) *DocumentService {
	return &DocumentService{DB: db}
}

type CreateDocumentInput struct {
	Title      string   `json:"title" binding:"required,max=256"`
	Content    string   `json:"content"`
	Summary    string   `json:"summary"`
	CategoryID *uint64  `json:"category_id"`
	TagIDs     []uint64 `json:"tag_ids"`
	Status     string   `json:"status"`
}

type UpdateDocumentInput struct {
	Title      string   `json:"title" binding:"required,max=256"`
	Content    string   `json:"content"`
	Summary    string   `json:"summary"`
	CategoryID *uint64  `json:"category_id"`
	TagIDs     []uint64 `json:"tag_ids"`
	Status     string   `json:"status"`
	IsFavorite *bool    `json:"is_favorite"`
}

// List 获取文档列表
func (s *DocumentService) List(userID uint64, page, pageSize int, categoryID *uint64, status string, tagID *uint64) (*utils.PaginatedData, error) {
	var total int64
	var documents []*models.Document

	query := s.DB.Model(&models.Document{}).Where("user_id = ?", userID)
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if tagID != nil {
		query = query.Joins("JOIN document_tags ON document_tags.document_id = documents.id").
			Where("document_tags.tag_id = ?", *tagID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("查询文档数量失败")
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Category").Preload("Tags").
		Order("updated_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&documents).Error
	if err != nil {
		return nil, errors.New("查询文档列表失败")
	}

	return &utils.PaginatedData{
		List: documents,
		Pagination: utils.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}, nil
}

// Search 全文搜索
func (s *DocumentService) Search(userID uint64, keyword string, page, pageSize int, tagID *uint64, categoryID *uint64) (*utils.PaginatedData, error) {
	var total int64
	var documents []*models.Document

	query := s.DB.Model(&models.Document{}).
		Where("user_id = ?", userID).
		Where(
			"title ILIKE ? OR content ILIKE ?",
			fmt.Sprintf("%%%s%%", keyword),
			fmt.Sprintf("%%%s%%", keyword),
		)
	if tagID != nil {
		query = query.Joins("JOIN document_tags ON document_tags.document_id = documents.id").
			Where("document_tags.tag_id = ?", *tagID)
	}
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("搜索失败")
	}

	err := query.Preload("Category").Preload("Tags").
		Order("updated_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&documents).Error
	if err != nil {
		return nil, errors.New("搜索失败")
	}

	return &utils.PaginatedData{
		List: documents,
		Pagination: utils.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}, nil
}

// GetByID 获取文档详情
func (s *DocumentService) GetByID(userID uint64, id uint64) (*models.Document, error) {
	var document models.Document
	if err := s.DB.Where("id = ? AND user_id = ?", id, userID).
		Preload("Category").Preload("Tags").
		First(&document).Error; err != nil {
		return nil, errors.New("文档不存在")
	}

	// 增加阅读次数
	s.DB.Model(&document).UpdateColumn("view_count", gorm.Expr("view_count + 1"))

	return &document, nil
}

// Create 创建文档
func (s *DocumentService) Create(userID uint64, input CreateDocumentInput) (*models.Document, error) {
	document := &models.Document{
		Title:      input.Title,
		Content:    input.Content,
		Summary:    input.Summary,
		CategoryID: input.CategoryID,
		UserID:     userID,
		Version:    1,
	}
	if input.Status != "" {
		document.Status = input.Status
	}

	tx := s.DB.Begin()

	if err := tx.Create(document).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("创建文档失败")
	}

	// 处理标签
	if len(input.TagIDs) > 0 {
		var tags []*models.Tag
		for _, tagID := range input.TagIDs {
			tags = append(tags, &models.Tag{ID: tagID})
		}
		if err := tx.Model(document).Association("Tags").Replace(tags); err != nil {
			tx.Rollback()
			return nil, errors.New("关联标签失败")
		}
	}

	// 创建初始版本
	version := &models.DocumentVersion{
		DocumentID: document.ID,
		Version:    1,
		Title:      input.Title,
		Content:    input.Content,
	}
	if err := tx.Create(version).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("创建版本记录失败")
	}

	tx.Commit()

	// 重新加载（含关联）
	s.DB.Where("id = ?", document.ID).Preload("Category").Preload("Tags").First(document)
	return document, nil
}

// Update 更新文档
func (s *DocumentService) Update(userID uint64, id uint64, input UpdateDocumentInput) (*models.Document, error) {
	var document models.Document
	if err := s.DB.Where("id = ? AND user_id = ?", id, userID).First(&document).Error; err != nil {
		return nil, errors.New("文档不存在")
	}

	newVersion := document.Version + 1

	updates := map[string]interface{}{
		"title":    input.Title,
		"content":  input.Content,
		"summary":  input.Summary,
		"category_id": input.CategoryID,
		"version":  newVersion,
	}
	if input.Status != "" {
		updates["status"] = input.Status
	}
	if input.IsFavorite != nil {
		updates["is_favorite"] = *input.IsFavorite
	}

	tx := s.DB.Begin()

	if err := tx.Model(&document).Updates(updates).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("更新文档失败")
	}

	// 更新标签
	if input.TagIDs != nil {
		var tags []*models.Tag
		for _, tagID := range input.TagIDs {
			tags = append(tags, &models.Tag{ID: tagID})
		}
		if err := tx.Model(&document).Association("Tags").Replace(tags); err != nil {
			tx.Rollback()
			return nil, errors.New("更新标签失败")
		}
	}

	// 保存版本记录
	version := &models.DocumentVersion{
		DocumentID: document.ID,
		Version:    newVersion,
		Title:      input.Title,
		Content:    input.Content,
	}
	if err := tx.Create(version).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("保存版本记录失败")
	}

	tx.Commit()

	s.DB.Where("id = ?", document.ID).Preload("Category").Preload("Tags").First(&document)
	return &document, nil
}

// Delete 删除文档
func (s *DocumentService) Delete(userID uint64, id uint64) error {
	result := s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Document{})
	if result.RowsAffected == 0 {
		return errors.New("文档不存在")
	}
	return result.Error
}

// GetVersions 获取文档版本列表
func (s *DocumentService) GetVersions(userID uint64, documentID uint64) ([]*models.DocumentVersion, error) {
	// 验证文档属于当前用户
	var count int64
	if err := s.DB.Model(&models.Document{}).Where("id = ? AND user_id = ?", documentID, userID).Count(&count).Error; err != nil || count == 0 {
		return nil, errors.New("文档不存在")
	}

	var versions []*models.DocumentVersion
	if err := s.DB.Where("document_id = ?", documentID).
		Order("version DESC").
		Find(&versions).Error; err != nil {
		return nil, errors.New("获取版本列表失败")
	}
	return versions, nil
}

// GetVersion 获取文档指定版本
func (s *DocumentService) GetVersion(userID uint64, documentID uint64, versionID uint64) (*models.DocumentVersion, error) {
	// 验证文档属于当前用户
	var count int64
	if err := s.DB.Model(&models.Document{}).Where("id = ? AND user_id = ?", documentID, userID).Count(&count).Error; err != nil || count == 0 {
		return nil, errors.New("文档不存在")
	}

	var version models.DocumentVersion
	if err := s.DB.Where("id = ? AND document_id = ?", versionID, documentID).First(&version).Error; err != nil {
		return nil, errors.New("版本不存在")
	}
	return &version, nil
}
