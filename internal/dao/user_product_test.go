package dao

import (
	"tamago/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *DAOSuite) TestGetUserProductsByProductID() {
	var (
		productID           uint64 = 1
		userProductID       uint64 = 1
		userID              uint64 = 1
		title                      = "first youtube plan"
		userProductMemberID uint64 = 1
	)

	userProductRows := sqlmock.NewRows([]string{"id", "host_id", "product_id", "title"}).
		AddRow(userProductID, userID, productID, title)
	s.Mock.ExpectQuery("SELECT").WithArgs(productID).WillReturnRows(userProductRows)

	userProductMemberRows := sqlmock.NewRows([]string{"id", "user_product_id"}).AddRow(userProductMemberID, userProductID)
	s.Mock.ExpectQuery("SELECT").WithArgs(userProductID).WillReturnRows(userProductMemberRows)

	userProducts, err := s.DAO.GetUserProductsByProductID(productID)

	require.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(*userProducts))
	expectedUserProduct := &models.UserProduct{}
	expectedUserProduct.ID = userProductID
	expectedUserProduct.ProductID = productID
	expectedUserProduct.HostID = userID
	expectedUserProduct.Title = title
	assert.Equal(s.T(), expectedUserProduct, &(*userProducts)[0])
}

func (s *DAOSuite) TestGetUserProductByUserProductID() {
	var (
		productID           uint64 = 1
		userProductID       uint64 = 1
		userID              uint64 = 1
		title                      = "first youtube plan"
		userProductMemberID uint64 = 1
		billingRequestID    uint64 = 1
	)

	s.Mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.
		NewRows([]string{"id", "host_id", "product_id", "title"}).
		AddRow(userProductID, userID, productID, title))
	s.Mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.
		NewRows([]string{"id", "user_id", "user_product_id"}).
		AddRow(userProductMemberID, userID, userProductID))
	s.Mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.
		NewRows([]string{"id", "user_product_id", "user_product_member_id"}).
		AddRow(billingRequestID, userProductID, userProductMemberID))

	up, err := s.DAO.GetUserProductByUserProductID(userProductID)

	require.Nil(s.T(), err)
	expectedUserProduct := &models.UserProduct{}
	expectedUserProduct.ID = userProductID
	expectedUserProduct.HostID = userID
	expectedUserProduct.ProductID = productID
	expectedUserProduct.Title = title
	expectedUserProductMember := &models.UserProductMember{}
	expectedUserProductMember.ID = userProductMemberID
	expectedUserProductMember.UserID = userID
	expectedUserProductMember.UserProductID = userProductID
	expectedUserProduct.UserProductMembers = []models.UserProductMember{*expectedUserProductMember}
	expectedBillingRequest := &models.BillingRequest{}
	expectedBillingRequest.ID = billingRequestID
	expectedBillingRequest.UserProductID = userProductID
	expectedBillingRequest.UserProductMemberID = userProductMemberID
	expectedUserProduct.BillingRequests = []models.BillingRequest{*expectedBillingRequest}

	assert.Equal(s.T(), expectedUserProduct, up)
}

func (s *DAOSuite) TestGetUserProductsByUserID() {
	var (
		userProductID       uint64 = 1
		hostID              uint64 = 1
		productID           uint64 = 1
		userProductMemberID uint64 = 1
		title                      = "awesome plan"
	)

	s.Mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.
		NewRows([]string{"id", "host_id", "product_id", "title"}).
		AddRow(userProductID, hostID, productID, title))
	s.Mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.
		NewRows([]string{"id", "user_id", "user_product_id"}).
		AddRow(userProductMemberID, hostID, userProductID))

	userProductMemeber := models.UserProductMember{
		UserID:        hostID,
		UserProductID: userProductID,
	}
	userProductMemeber.ID = userProductMemberID
	userProduct := models.UserProduct{
		HostID:             hostID,
		ProductID:          productID,
		Title:              title,
		UserProductMembers: []models.UserProductMember{userProductMemeber},
	}
	userProduct.ID = userProductID
	expectedUserProducts := []models.UserProduct{userProduct}

	userProducts, err := s.DAO.GetUserProductsByUserID(hostID)

	require.Nil(s.T(), err)
	assert.Equal(s.T(), expectedUserProducts, userProducts)
}
