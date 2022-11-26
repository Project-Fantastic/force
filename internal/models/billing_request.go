package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type BillingRequest struct {
	BaseModel
	UserProductID       uint64  `gorm:"not null"`
	UserProductMemberID uint64  `gorm:"not null"`
	Status              uint    `gorm:"not null;default:'0'"`
	Amount              float64 `gorm:"not null"`
	StartDate           *time.Time
	EndDate             *time.Time
	ExpirationDate      *time.Time
}

func (br *BillingRequest) GetBillingRequestByUserProductID(db *gorm.DB, upID uint64) (*[]BillingRequest, error) {
	var brs []BillingRequest
	err := db.Where("user_product_id = ?", upID).First(&brs).Error
	return &brs, err
}
