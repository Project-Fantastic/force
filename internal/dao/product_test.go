package dao

import (
	"fmt"
	"tamago/internal/api"
	"tamago/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *DAOSuite) TestGetProductByID() {
	var (
		productID uint64 = 1
		name             = "youtube"
	)

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(productID, name)

	s.Mock.ExpectQuery("SELECT").WillReturnRows(rows)

	product := &models.Product{}
	product, err := product.GetProductByProductID(s.DB, productID)

	expectedProduct := &models.Product{}
	expectedProduct.ID = productID
	expectedProduct.Name = name

	require.Nil(s.T(), err)
	assert.EqualValues(s.T(), product, expectedProduct)
}

func (s *DAOSuite) TestGetProductsByBillingTypeWithMultipleResults() {
	var (
		recurringType = api.Product_BILLING_TYPE_RECURRING
	)

	rows1 := sqlmock.NewRows([]string{"id", "name", "billing_type"}).
		AddRow(1, "youtube", recurringType).
		AddRow(2, "spotify", recurringType)

	s.Mock.ExpectQuery("^SELECT").WithArgs(recurringType).WillReturnRows(rows1)
	p := &models.Product{}
	products, err := p.GetProductsByBillingType(s.DB, recurringType)

	require.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(*products))
}

func (s *DAOSuite) TestGetProductsByBillingTypeWithSingleResult() {
	var (
		oneTimeType = api.Product_BILLING_TYPE_ONE_TIME
	)

	rows3 := sqlmock.NewRows([]string{"id", "name", "billing_type"}).AddRow(3, "turbo tax", oneTimeType)

	p := &models.Product{}
	s.Mock.ExpectQuery("SELECT").WithArgs(oneTimeType).WillReturnRows(rows3)
	products2, _ := p.GetProductsByBillingType(s.DB, oneTimeType)

	expectedProduct := &models.Product{}
	expectedProduct.ID = 3
	expectedProduct.Name = "turbo tax"
	expectedProduct.BillingType = int32(oneTimeType)

	assert.Equal(s.T(), 1, len(*products2))
	assert.Equal(s.T(), expectedProduct, &(*products2)[0])
}

func (s *DAOSuite) TestGetProductByNameExisting() {
	rows1 := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "youtube")

	p := &models.Product{}
	s.Mock.ExpectQuery("SELECT").WillReturnRows(rows1)
	p, err := p.GetProductByName(s.DB, "youtube")

	fmt.Printf("%+v", p)
	require.Nil(s.T(), err)
	assert.Equal(s.T(), uint64(1), p.ID)
	assert.Equal(s.T(), "youtube", p.Name)
}
