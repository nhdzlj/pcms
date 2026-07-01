package models

import (
	"time"
)

type Category struct {
	ID        uint64     `gorm:"primaryKey" json:"id"`
	Name      string     `gorm:"size:128;not null" json:"name"`
	ParentID  *uint64    `gorm:"index" json:"parent_id"`
	SortOrder int        `gorm:"default:0" json:"sort_order"`
	UserID    uint64     `gorm:"index;not null" json:"user_id"`
	Icon      string     `gorm:"size:64;default:folder" json:"icon"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// 关联
	Children  []*Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Parent    *Category   `gorm:"foreignKey:ParentID" json:"-"`
	User      *User       `gorm:"foreignKey:UserID" json:"-"`
}

func (Category) TableName() string {
	return "categories"
}

// BuildTree 将扁平列表构建为树形结构
func BuildTree(categories []*Category) []*Category {
	nodeMap := make(map[uint64]*Category)
	var roots []*Category

	for _, c := range categories {
		c.Children = make([]*Category, 0)
		nodeMap[c.ID] = c
	}

	for _, c := range categories {
		if c.ParentID == nil {
			roots = append(roots, c)
		} else if parent, ok := nodeMap[*c.ParentID]; ok {
			parent.Children = append(parent.Children, c)
		}
	}

	return roots
}
