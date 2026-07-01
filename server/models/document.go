package models

import (
	"time"
)

type Document struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	Title      string    `gorm:"size:256;not null" json:"title"`
	Content    string    `gorm:"type:text" json:"content"`
	Summary    string    `gorm:"type:text" json:"summary"`
	CategoryID *uint64   `gorm:"index" json:"category_id"`
	UserID     uint64    `gorm:"index;not null" json:"user_id"`
	Status     string    `gorm:"size:20;default:draft" json:"status"`
	ViewCount  int       `gorm:"default:0" json:"view_count"`
	IsFavorite bool      `gorm:"default:false" json:"is_favorite"`
	Version    int       `gorm:"default:1" json:"version"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// 关联
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags     []*Tag    `gorm:"many2many:document_tags;" json:"tags,omitempty"`
	User     *User     `gorm:"foreignKey:UserID" json:"-"`
}

func (Document) TableName() string {
	return "documents"
}
