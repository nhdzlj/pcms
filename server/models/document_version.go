package models

import "time"

type DocumentVersion struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	DocumentID uint64    `gorm:"index;not null" json:"document_id"`
	Version    int       `gorm:"not null" json:"version"`
	Title      string    `gorm:"size:256" json:"title"`
	Content    string    `gorm:"type:text" json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

func (DocumentVersion) TableName() string {
	return "document_versions"
}
