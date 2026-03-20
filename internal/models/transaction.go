package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID            string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID        string         `gorm:"type:uuid;not null;index" json:"userId"`
	User          User           `gorm:"foreignKey:UserID" json:"-"`
	CategoryID    string         `gorm:"type:uuid;not null" json:"categoryId"`
	Category      Category       `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"category"`
	LinkedGoalID  *string        `gorm:"type:uuid" json:"linkedGoalId"`
	Title         string         `gorm:"type:varchar(255);not null" json:"title"`
	Amount        float64        `gorm:"type:decimal(15,2)" json:"amount"`
	Type          string         `gorm:"type:varchar(20);not null" json:"type"` // income, expense, transfer
	Date          time.Time      `gorm:"not null;index" json:"date"`
	PaymentMethod string         `gorm:"type:varchar(50)" json:"paymentMethod"`
	Note          string         `gorm:"type:text" json:"note"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
