package services

import (
	"errors"

	"pcms/models"
	"pcms/store"

	"gorm.io/gorm"
)

type CategoryService struct {
	DB store.Store
}

func NewCategoryService(db store.Store) *CategoryService {
	return &CategoryService{DB: db}
}

type CreateCategoryInput struct {
	Name     string  `json:"name" binding:"required,max=128"`
	ParentID *uint64 `json:"parent_id"`
	Icon     string  `json:"icon"`
}

type UpdateCategoryInput struct {
	Name      string  `json:"name" binding:"required,max=128"`
	ParentID  *uint64 `json:"parent_id"`
	SortOrder *int    `json:"sort_order"`
	Icon      string  `json:"icon"`
}

type MoveCategoryInput struct {
	ParentID  *uint64 `json:"parent_id"`
	SortOrder int     `json:"sort_order"`
}

func (s *CategoryService) GetTree(userID uint64) ([]*models.Category, error) {
	var categories []*models.Category
	if err := s.DB.Where("UserID = ?", userID).
		Order("SortOrder ASC, ID ASC").
		Find(&categories); err != nil {
		return nil, errors.New("获取分类失败")
	}
	return models.BuildTree(categories), nil
}

func (s *CategoryService) Create(userID uint64, input CreateCategoryInput) (*models.Category, error) {
	if input.ParentID != nil {
		var parent models.Category
		if err := s.DB.Where("ID = ? AND UserID = ?", *input.ParentID, userID).First(&parent); err != nil {
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

	if err := s.DB.Create(category); err != nil {
		return nil, errors.New("创建分类失败")
	}
	return category, nil
}

func (s *CategoryService) Update(userID uint64, id uint64, input UpdateCategoryInput) (*models.Category, error) {
	var category models.Category
	if err := s.DB.Where("ID = ? AND UserID = ?", id, userID).First(&category); err != nil {
		return nil, errors.New("分类不存在")
	}

	updates := map[string]interface{}{
		"Name": input.Name,
	}
	if input.ParentID != nil {
		if *input.ParentID == id {
			return nil, errors.New("不能将分类移动到自身下")
		}
		updates["ParentID"] = *input.ParentID
	}
	if input.SortOrder != nil {
		updates["SortOrder"] = *input.SortOrder
	}
	if input.Icon != "" {
		updates["Icon"] = input.Icon
	}

	if err := s.DB.Model(&models.Category{}).Where("ID = ? AND UserID = ?", id, userID).Updates(updates); err != nil {
		return nil, errors.New("更新分类失败")
	}

	s.DB.First(&category, id)
	return &category, nil
}

func (s *CategoryService) Move(userID uint64, id uint64, input MoveCategoryInput) (*models.Category, error) {
	var category models.Category
	if err := s.DB.Where("ID = ? AND UserID = ?", id, userID).First(&category); err != nil {
		return nil, errors.New("分类不存在")
	}

	if input.ParentID != nil && *input.ParentID == id {
		return nil, errors.New("不能将分类移动到自身下")
	}

	updates := map[string]interface{}{}
	if input.ParentID != nil {
		updates["ParentID"] = *input.ParentID
	}
	updates["SortOrder"] = input.SortOrder

	if err := s.DB.Model(&models.Category{}).Where("ID = ? AND UserID = ?", id, userID).Updates(updates); err != nil {
		return nil, errors.New("移动分类失败")
	}

	s.DB.First(&category, id)
	return &category, nil
}

func (s *CategoryService) Delete(userID uint64, id uint64) error {
	var category models.Category
	if err := s.DB.Where("ID = ? AND UserID = ?", id, userID).First(&category); err != nil {
		return errors.New("分类不存在")
	}

	if err := s.DB.Delete(&category); err != nil {
		return errors.New("删除分类失败")
	}
	return nil
}

// 防止未使用导入（gorm 在 services 包中仍被 document.go 需要）
var _ = gorm.ErrRecordNotFound
