package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type UserProductMember struct {
	BaseModel
	UserID        uint64 `gorm:"not null"`
	UserProductID uint64 `gorm:"not null"`
	IsHost        bool   `gorm:"not null;default: 'false'"`
	//Status:
	// 0: pending request
	// 1: accepted/delivered
	// 2: rejected
	// 3: deliver_confirmed
	Status          uint `gorm:"not null;default: '0'"`
	BillingRequests []BillingRequest
}

func (m *UserProductMember) prepare() {
	currentTime := time.Now().UTC()
	m.ID = 0
	m.CreatedAt = currentTime
	m.UpdatedAt = currentTime
	m.Status = 0
}

//GetUserProductMembersByUserProductID returns a slice of members for a specific user product
func (m *UserProductMember) GetUserProductMembersByUserProductID(
	db *gorm.DB,
	upID uint64) (*[]UserProductMember, error) {
	var members []UserProductMember
	err := db.Where("user_product_id = ?", upID).Find(&members).Error

	return &members, err
}

//GetUserProductMemberByID returns a user product member by ID
func (m *UserProductMember) GetUserProductMemberByID(db *gorm.DB, upmID uint64) (*UserProductMember, error) {
	err := db.First(m, upmID).Error
	return m, err
}

//CreateUserProductMember creates a new user product member and returns it
func (m *UserProductMember) CreateUserProductMember(db *gorm.DB) (*UserProductMember, error) {
	m.prepare()
	err := db.Create(m).Error
	return m, err
}

//UpdateUserProductMember updates an existing user product member, such as when the status changes
func (m *UserProductMember) UpdateUserProductMember(db *gorm.DB) (*UserProductMember, error) {
	err := db.Save(m).Error
	return m, err
}

//DeleteUserProductMember deletes an existing user product member by adding a delete_at field
func (m *UserProductMember) DeleteUserProductMember(db *gorm.DB) error {
	return db.Delete(m).Error
}
