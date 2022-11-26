package models

type UserThirdPartyPayment struct {
	BaseModel
	UserID        uint64 `gorm:"not null"`
	PaymentType   uint   `gorm:"not null;default:'0'"`
	PaymentID     string `gorm:"type:varchar(512);not null"`
	PaymentIDType uint   `gorm:"not null;default:'0'"`
}
