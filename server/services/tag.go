package services

import (
	"errors"

	"pcms/models"
	"pcms/store"
)

type TagService struct {
	DB store.Store
}

func NewTagService(db store.Store) *TagService {
	return &TagService{DB: db}
}

type CreateTagInput struct {
	Name string `json:"name" binding:"required,max=64"`
}

func (s *TagService) List(userID uint64) ([]*models.Tag, error) {
	var tags []*models.Tag
	if err := s.DB.Where("UserID = ?", userID).Order("Name ASC").Find(&tags); err != nil {
		return nil, errors.New("获取标签失败")
	}
	return tags, nil
}

func (s *TagService) Create(userID uint64, input CreateTagInput) (*models.Tag, error) {
	// FirstOrCreate: 先查找
	var existing models.Tag
	if err := s.DB.Where("Name = ? AND UserID = ?", input.Name, userID).First(&existing); err == nil {
		return &existing, nil
	}

	tag := &models.Tag{
		Name:   input.Name,
		UserID: userID,
	}
	if err := s.DB.Create(tag); err != nil {
		return nil, errors.New("创建标签失败")
	}
	return tag, nil
}

func (s *TagService) Delete(userID uint64, id uint64) error {
	if err := s.DB.Delete(&models.Tag{}, id, "UserID = ?", userID); err != nil {
		return errors.New("标签不存在")
	}
	return nil
}
