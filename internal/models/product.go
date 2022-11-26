package models

import (
	"tamago/internal/api"

	"github.com/jinzhu/gorm"
)

type Product struct {
	BaseModel
	Name                string `gorm:"type:varchar(512);not null;unique_index"`
	BillingType         int32  `gorm:"not null;default:'0'"`
	IsFixedPrice        bool   `gorm:"not null;default:'true'"`
	MaxMemberCount      uint64 `gorm:"not null;default:'1'"`
	AccountSupportTypes uint   `gorm:"default:'0'"`
	UserProducts        []UserProduct
	UserProductAccounts []UserProductAccount
}

//GetProductByProductID returns a single Product based on product ID
func (p *Product) GetProductByProductID(db *gorm.DB, id uint64) (*Product, error) {
	err := db.First(p, id).Error
	return p, err
}

//GetProductsByBillingType returns a slice of Products based on billingType: one_time or recurring
func (p *Product) GetProductsByBillingType(db *gorm.DB, billingType api.Product_BillingType) (*[]Product, error) {
	var products []Product
	err := db.Where("billing_type = ?", billingType).Find(&products).Error
	return &products, err
}

//GetProductByName returns a slice of Products based on name, used for search
func (p *Product) GetProductByName(db *gorm.DB, name string) (*Product, error) {
	err := db.Where("name LIKE \"%?%\"", name).First(p).Error
	return p, err
}
