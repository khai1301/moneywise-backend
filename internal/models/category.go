package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID      string         `gorm:"type:uuid;not null;index" json:"userId"`
	User        User           `gorm:"foreignKey:UserID" json:"-"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Type        string         `gorm:"type:varchar(20);not null" json:"type"` // income, expense, both
	Icon        string         `gorm:"type:varchar(255)" json:"icon"`
	Color       string         `gorm:"type:varchar(50)" json:"color"`
	Description string         `gorm:"type:text" json:"description"`
	IsSystem    bool           `gorm:"default:false" json:"isSystem"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
