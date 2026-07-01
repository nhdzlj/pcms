package services

import (
	"errors"

	"pcms/models"

	"gorm.io/gorm"
)

type TagService struct {
	DB *gorm.DB
}

func NewTagService(db *gorm.DB) *TagService {
	return &TagService{DB: db}
}

type CreateTagInput struct {
	Name string `json:"name" binding:"required,max=64"`
}

// List 获取标签列表
func (s *TagService) List(userID uint64) ([]*models.Tag, error) {
	var tags []*models.Tag
	if err := s.DB.Where("user_id = ?", userID).Order("name ASC").Find(&tags).Error; err != nil {
		return nil, errors.New("获取标签失败")
	}
	return tags, nil
}

// Create 创建标签
func (s *TagService) Create(userID uint64, input CreateTagInput) (*models.Tag, error) {
	tag := &models.Tag{
		Name:   input.Name,
		UserID: userID,
	}
	if err := s.DB.Where("name = ? AND user_id = ?", input.Name, userID).
		FirstOrCreate(tag).Error; err != nil {
		return nil, errors.New("创建标签失败")
	}
	return tag, nil
}

// Delete 删除标签
func (s *TagService) Delete(userID uint64, id uint64) error {
	result := s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Tag{})
	if result.RowsAffected == 0 {
		return errors.New("标签不存在")
	}
	return result.Error
}
