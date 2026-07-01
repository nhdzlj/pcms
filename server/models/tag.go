package models

type Tag struct {
	ID     uint64 `gorm:"primaryKey" json:"id"`
	Name   string `gorm:"size:64;not null" json:"name"`
	UserID uint64 `gorm:"index;not null" json:"user_id"`
}

func (Tag) TableName() string {
	return "tags"
}
