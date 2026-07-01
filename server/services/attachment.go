package services

import (
	"errors"

	"pcms/models"

	"gorm.io/gorm"
)

type AttachmentService struct {
	DB *gorm.DB
}

func NewAttachmentService(db *gorm.DB) *AttachmentService {
	return &AttachmentService{DB: db}
}

type CreateAttachmentInput struct {
	FileName string `json:"file_name" binding:"required"`
	FilePath string `json:"file_path" binding:"required"`
	FileSize int64  `json:"file_size"`
	MimeType string `json:"mime_type"`
}

type ListAttachmentQuery struct {
	DocumentID *uint64
}

// List 获取附件列表
func (s *AttachmentService) List(userID uint64, query ListAttachmentQuery) ([]*models.Attachment, error) {
	var attachments []*models.Attachment
	q := s.DB.Where("user_id = ?", userID)
	if query.DocumentID != nil {
		q = q.Where("document_id = ?", *query.DocumentID)
	}
	if err := q.Order("created_at DESC").Find(&attachments).Error; err != nil {
		return nil, errors.New("获取附件列表失败")
	}
	return attachments, nil
}

// Create 创建附件记录
func (s *AttachmentService) Create(userID uint64, input CreateAttachmentInput, documentID *uint64) (*models.Attachment, error) {
	attachment := &models.Attachment{
		DocumentID: documentID,
		FileName:   input.FileName,
		FilePath:   input.FilePath,
		FileSize:   input.FileSize,
		MimeType:   input.MimeType,
		UserID:     userID,
	}
	if err := s.DB.Create(attachment).Error; err != nil {
		return nil, errors.New("创建附件记录失败")
	}
	return attachment, nil
}

// Delete 删除附件
func (s *AttachmentService) Delete(userID uint64, id uint64) error {
	result := s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Attachment{})
	if result.RowsAffected == 0 {
		return errors.New("附件不存在")
	}
	return result.Error
}

// BindDocument 关联附件到文档
func (s *AttachmentService) BindDocument(userID uint64, id uint64, documentID uint64) error {
	result := s.DB.Model(&models.Attachment{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("document_id", documentID)
	if result.RowsAffected == 0 {
		return errors.New("附件不存在")
	}
	return result.Error
}
