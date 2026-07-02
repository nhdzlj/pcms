package services

import (
	"errors"

	"pcms/models"
	"pcms/store"
)

type AttachmentService struct {
	DB store.Store
}

func NewAttachmentService(db store.Store) *AttachmentService {
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

func (s *AttachmentService) List(userID uint64, query ListAttachmentQuery) ([]*models.Attachment, error) {
	var attachments []*models.Attachment
	q := s.DB.Where("UserID = ?", userID)
	if query.DocumentID != nil {
		q = q.Where("DocumentID = ?", *query.DocumentID)
	}
	if err := q.Order("CreatedAt DESC").Find(&attachments); err != nil {
		return nil, errors.New("获取附件列表失败")
	}
	return attachments, nil
}

func (s *AttachmentService) Create(userID uint64, input CreateAttachmentInput, documentID *uint64) (*models.Attachment, error) {
	attachment := &models.Attachment{
		DocumentID: documentID,
		FileName:   input.FileName,
		FilePath:   input.FilePath,
		FileSize:   input.FileSize,
		MimeType:   input.MimeType,
		UserID:     userID,
	}
	if err := s.DB.Create(attachment); err != nil {
		return nil, errors.New("创建附件记录失败")
	}
	return attachment, nil
}

func (s *AttachmentService) Delete(userID uint64, id uint64) error {
	if err := s.DB.Delete(&models.Attachment{}, id, "UserID = ?", userID); err != nil {
		return errors.New("附件不存在")
	}
	return nil
}

func (s *AttachmentService) BindDocument(userID uint64, id uint64, documentID uint64) error {
	if err := s.DB.Model(&models.Attachment{}).Where("ID = ? AND UserID = ?", id, userID).
		Updates(map[string]interface{}{"DocumentID": documentID}); err != nil {
		return errors.New("附件不存在")
	}
	return nil
}
