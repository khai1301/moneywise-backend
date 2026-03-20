package models

import (
	"time"

	"gorm.io/gorm"
)

// Budget defines a spending limit for a specific category in a given month.
// Month is stored as "YYYY-MM" string for simplicity.
type Budget struct {
	ID         string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID     string         `gorm:"type:uuid;not null;index" json:"userId"`
	User       User           `gorm:"foreignKey:UserID" json:"-"`
	CategoryID string         `gorm:"type:uuid;not null" json:"categoryId"`
	Category   Category       `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"category"`
	Amount     float64        `gorm:"type:decimal(15,2);not null" json:"amount"` // Budget limit
	Month      string         `gorm:"type:varchar(7);not null;index" json:"month"` // "YYYY-MM"
	Note       string         `gorm:"type:text" json:"note"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
