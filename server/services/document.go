package services

import (
	"errors"
	"fmt"

	"pcms/models"
	"pcms/store"
	"pcms/utils"
)

type DocumentService struct {
	DB store.Store
}

func NewDocumentService(db store.Store) *DocumentService {
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

func (s *DocumentService) List(userID uint64, page, pageSize int, categoryID *uint64, status string, tagID *uint64) (*utils.PaginatedData, error) {
	var total int64
	var documents []*models.Document

	query := s.DB.Model(&models.Document{}).Where("UserID = ?", userID)
	if categoryID != nil {
		query = query.Where("CategoryID = ?", *categoryID)
	}
	if status != "" {
		query = query.Where("Status = ?", status)
	}
	if tagID != nil {
		query = query.Joins("JOIN document_tags ON document_tags.document_id = documents.id").
			Where("document_tags.tag_id = ?", *tagID)
	}

	if err := query.Count(&total); err != nil {
		return nil, errors.New("查询文档数量失败")
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("Category").Preload("Tags").
		Order("UpdatedAt DESC").
		Offset(offset).Limit(pageSize).
		Find(&documents); err != nil {
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

func (s *DocumentService) Search(userID uint64, keyword string, page, pageSize int, tagID *uint64, categoryID *uint64) (*utils.PaginatedData, error) {
	var total int64
	var documents []*models.Document

	query := s.DB.Model(&models.Document{}).
		Where("UserID = ?", userID).
		Where(
			"Title LIKE ? OR Content LIKE ?",
			fmt.Sprintf("%%%s%%", keyword),
			fmt.Sprintf("%%%s%%", keyword),
		)
	if tagID != nil {
		query = query.Joins("JOIN document_tags ON document_tags.document_id = documents.id").
			Where("document_tags.tag_id = ?", *tagID)
	}
	if categoryID != nil {
		query = query.Where("CategoryID = ?", *categoryID)
	}

	if err := query.Count(&total); err != nil {
		return nil, errors.New("搜索失败")
	}

	err := query.Preload("Category").Preload("Tags").
		Order("UpdatedAt DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&documents)
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

func (s *DocumentService) GetByID(userID uint64, id uint64) (*models.Document, error) {
	var document models.Document
	if err := s.DB.Where("ID = ? AND UserID = ?", id, userID).
		Preload("Category").Preload("Tags").
		First(&document); err != nil {
		return nil, errors.New("文档不存在")
	}

	s.DB.Model(&models.Document{}).Where("ID = ?", id).UpdateColumn("ViewCount", "ViewCount + 1")

	return &document, nil
}

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

	if err := s.DB.Create(document); err != nil {
		return nil, errors.New("创建文档失败")
	}

	// 重新加载（含关联）
	var doc models.Document
	s.DB.Where("ID = ?", document.ID).Preload("Category").Preload("Tags").First(&doc)
	return &doc, nil
}

func (s *DocumentService) Update(userID uint64, id uint64, input UpdateDocumentInput) (*models.Document, error) {
	var document models.Document
	if err := s.DB.Where("ID = ? AND UserID = ?", id, userID).First(&document); err != nil {
		return nil, errors.New("文档不存在")
	}

	newVersion := document.Version + 1

	updates := map[string]interface{}{
		"Title":      input.Title,
		"Content":    input.Content,
		"Summary":    input.Summary,
		"CategoryID": input.CategoryID,
		"Version":    newVersion,
	}
	if input.Status != "" {
		updates["Status"] = input.Status
	}
	if input.IsFavorite != nil {
		updates["IsFavorite"] = *input.IsFavorite
	}

	if err := s.DB.Model(&models.Document{}).Where("ID = ? AND UserID = ?", id, userID).Updates(updates); err != nil {
		return nil, errors.New("更新文档失败")
	}

	var doc models.Document
	s.DB.Where("ID = ?", id).Preload("Category").Preload("Tags").First(&doc)
	return &doc, nil
}

func (s *DocumentService) Delete(userID uint64, id uint64) error {
	if err := s.DB.Delete(&models.Document{}, id, "UserID = ?", userID); err != nil {
		return errors.New("文档不存在")
	}
	return nil
}

func (s *DocumentService) GetVersions(userID uint64, documentID uint64) ([]*models.DocumentVersion, error) {
	var count int64
	if err := s.DB.Model(&models.Document{}).Where("ID = ? AND UserID = ?", documentID, userID).Count(&count); err != nil || count == 0 {
		return nil, errors.New("文档不存在")
	}

	var versions []*models.DocumentVersion
	if err := s.DB.Where("DocumentID = ?", documentID).
		Order("Version DESC").
		Find(&versions); err != nil {
		return nil, errors.New("获取版本列表失败")
	}
	return versions, nil
}

func (s *DocumentService) GetVersion(userID uint64, documentID uint64, versionID uint64) (*models.DocumentVersion, error) {
	var count int64
	if err := s.DB.Model(&models.Document{}).Where("ID = ? AND UserID = ?", documentID, userID).Count(&count); err != nil || count == 0 {
		return nil, errors.New("文档不存在")
	}

	var version models.DocumentVersion
	if err := s.DB.Where("ID = ? AND DocumentID = ?", versionID, documentID).First(&version); err != nil {
		return nil, errors.New("版本不存在")
	}
	return &version, nil
}

// 防止未使用的导入
var _ = fmt.Sprintf
