package models

import "github.com/jinzhu/gorm"

type UserProduct struct {
	BaseModel
	HostID             uint64  `gorm:"not null"`
	ProductID          uint64  `gorm:"not null"`
	Title              string  `gorm:"type:varchar(512);not null"`
	Description        string  `gorm:"type:text"`
	Active             bool    `gorm:"not null;default:'true'"`
	MinMemberCount     uint64  `gorm:"not null;default:'1'"`
	MaxMemberCount     uint64  `gorm:"not null;default:'1'"`
	MinPrice           float64 `gorm:"not null"`
	MaxPrice           float64 `gorm:"not null"`
	TotalPrice         float64 `gorm:"not null"`
	UserProductMembers []UserProductMember
	BillingRequests    []BillingRequest
}

//GetUserProductsByProductID returns a slice of user products based on product id
func (u *UserProduct) GetUserProductsByProductID(db *gorm.DB, productID uint64) (*[]UserProduct, error) {
	var userProducts []UserProduct
	err := db.Where("product_id = ?", productID).Find(&userProducts).Error
	return &userProducts, err
}

//GetUserProductByUserProductID returns a user products by user product id
func (u *UserProduct) GetUserProductByUserProductID(db *gorm.DB, id uint64) (*UserProduct, error) {
	err := db.First(u, id).Error
	return u, err
}

// GetUserProductsByUserID returns a slice of user products that the user owns or joins
func (u *UserProduct) GetUserProductsByUserID(db *gorm.DB, userID uint64) ([]UserProduct, error) {
	var userProducts []UserProduct
	err := db.Preload("UserProductMembers", "user_id = ?", userID).Find(&userProducts).Error
	return userProducts, err
}
