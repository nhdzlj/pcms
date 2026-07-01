package services

import (
	"errors"

	"pcms/models"

	"gorm.io/gorm"
)

type CategoryService struct {
	DB *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{DB: db}
}

type CreateCategoryInput struct {
	Name     string  `json:"name" binding:"required,max=128"`
	ParentID *uint64 `json:"parent_id"`
	Icon     string  `json:"icon"`
}

type UpdateCategoryInput struct {
	Name     string  `json:"name" binding:"required,max=128"`
	ParentID *uint64 `json:"parent_id"`
	SortOrder *int   `json:"sort_order"`
	Icon     string  `json:"icon"`
}

type MoveCategoryInput struct {
	ParentID  *uint64 `json:"parent_id"`
	SortOrder int     `json:"sort_order"`
}

// GetTree 获取分类树
func (s *CategoryService) GetTree(userID uint64) ([]*models.Category, error) {
	var categories []*models.Category
	err := s.DB.Where("user_id = ?", userID).
		Order("sort_order ASC, id ASC").
		Find(&categories).Error
	if err != nil {
		return nil, errors.New("获取分类失败")
	}
	return models.BuildTree(categories), nil
}

// Create 创建分类
func (s *CategoryService) Create(userID uint64, input CreateCategoryInput) (*models.Category, error) {
	if input.ParentID != nil {
		// 验证父分类是否存在且属于当前用户
		var parent models.Category
		if err := s.DB.Where("id = ? AND user_id = ?", *input.ParentID, userID).First(&parent).Error; err != nil {
			return nil, errors.New("父分类不存在")
		}
	}

	category := &models.Category{
		Name:     input.Name,
		ParentID: input.ParentID,
		UserID:   userID,
		Icon:     input.Icon,
	}
	if category.Icon == "" {
		category.Icon = "folder"
	}

	if err := s.DB.Create(category).Error; err != nil {
		return nil, errors.New("创建分类失败")
	}

	return category, nil
}

// Update 更新分类
func (s *CategoryService) Update(userID uint64, id uint64, input UpdateCategoryInput) (*models.Category, error) {
	var category models.Category
	if err := s.DB.Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		return nil, errors.New("分类不存在")
	}

	updates := map[string]interface{}{
		"name": input.Name,
	}
	if input.ParentID != nil {
		// 不能将自己设为自己的子节点
		if *input.ParentID == id {
			return nil, errors.New("不能将分类移动到自身下")
		}
		updates["parent_id"] = *input.ParentID
	}
	if input.SortOrder != nil {
		updates["sort_order"] = *input.SortOrder
	}
	if input.Icon != "" {
		updates["icon"] = input.Icon
	}

	if err := s.DB.Model(&category).Updates(updates).Error; err != nil {
		return nil, errors.New("更新分类失败")
	}

	s.DB.First(&category, id)
	return &category, nil
}

// Move 移动分类
func (s *CategoryService) Move(userID uint64, id uint64, input MoveCategoryInput) (*models.Category, error) {
	var category models.Category
	if err := s.DB.Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		return nil, errors.New("分类不存在")
	}

	if input.ParentID != nil && *input.ParentID == id {
		return nil, errors.New("不能将分类移动到自身下")
	}

	updates := map[string]interface{}{}
	if input.ParentID != nil {
		updates["parent_id"] = *input.ParentID
	}
	updates["sort_order"] = input.SortOrder

	if err := s.DB.Model(&category).Updates(updates).Error; err != nil {
		return nil, errors.New("移动分类失败")
	}

	s.DB.First(&category, id)
	return &category, nil
}

// Delete 删除分类（及其子分类）
func (s *CategoryService) Delete(userID uint64, id uint64) error {
	var category models.Category
	if err := s.DB.Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		return errors.New("分类不存在")
	}

	if err := s.DB.Delete(&category).Error; err != nil {
		return errors.New("删除分类失败")
	}

	return nil
}
