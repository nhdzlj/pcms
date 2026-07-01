package models

import "time"

type Attachment struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	DocumentID *uint64   `gorm:"index" json:"document_id"`
	FileName   string    `gorm:"size:256" json:"file_name"`
	FilePath   string    `gorm:"size:512" json:"file_path"`
	FileSize   int64     `json:"file_size"`
	MimeType   string    `gorm:"size:128" json:"mime_type"`
	UserID     uint64    `gorm:"index;not null" json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (Attachment) TableName() string {
	return "attachments"
}
