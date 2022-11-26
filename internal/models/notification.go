package models

type Notification struct {
	BaseModel
	UserID           uint64 `gorm:"not null"`
	NotificationType uint   `gorm:"not null"`
	ObjectID         uint64 `gorm:"not null"`
	Status           uint   `gorm:"not null;default:'0'"`
}
