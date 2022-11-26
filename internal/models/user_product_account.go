package models

type UserProductAccount struct {
	BaseModel
	UserID    uint64 `gorm:"not null"`
	ProductID uint64 `gorm:"not null"`
	Account   string `gorm:"type:varchar(512);not null"`
}
