package dao

import (
	"tamago/internal/api"
	"tamago/internal/models"
)

//GetProductsByBillingType returns a slice of products by billingType
func (d *DAO) GetProductsByBillingType(billingType api.Product_BillingType) (*[]models.Product, error) {
	product := &models.Product{}
	return product.GetProductsByBillingType(d.db, billingType)
}

//GetProductByProductID returns a product by productId
func (d *DAO) GetProductByProductID(id uint64) (*models.Product, error) {
	product := &models.Product{}
	return product.GetProductByProductID(d.db, id)
}

//GetProductByName returns a slice of Products based on name, used for search
func (d *DAO) GetProductByName(name string) (*models.Product, error) {
	product := &models.Product{}
	return product.GetProductByName(d.db, name)
}
