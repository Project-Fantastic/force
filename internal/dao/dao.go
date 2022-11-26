package dao

import (
	"tamago/internal/api"
	"tamago/internal/models"

	"github.com/jinzhu/gorm"
)

// DAO is an object that talks with DB.
type DAO struct {
	db *gorm.DB
}

// DataAccessIface is an interface for DAO and defines all the DB calls.
type DataAccessIface interface {
	// User
	GetUserByID(uint64) (*models.User, error)
	GetUserByEmail(string) (*models.User, error)
	SignUpUser(string, string) (*models.User, error)
	VerifyLogin(string, string) (bool, uint64)

	// Product
	GetProductsByBillingType(api.Product_BillingType) (*[]models.Product, error)
	GetProductByProductID(uint64) (*models.Product, error)
	GetProductByName(string) (*models.Product, error)

	//UserProduct
	GetUserProductsByProductID(uint64) (*[]models.UserProduct, error)
	GetUserProductByUserProductID(uint64) (*models.UserProduct, error)
	GetUserProductsByUserID(uint64) ([]models.UserProduct, error)

	//Member
	GetUserProductMembersByUserProductID(uint64) (*[]models.UserProductMember, error)
	GetUserProductMemberByID(uint64) (*models.UserProductMember, error)
	CreateUserProductMember(uint64, uint64, bool) (*models.UserProductMember, error)
}

// NewDAO creates a new DAO object
func NewDAO(db *gorm.DB) *DAO {
	return &DAO{db: db}
}
