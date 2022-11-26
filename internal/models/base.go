package models

import "time"

type BaseModel struct {
	ID        uint64     `gorm:"primary_key"`
	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt time.Time  `gorm:"not null"`
	DeletedAt *time.Time `sql:"index"`
}
