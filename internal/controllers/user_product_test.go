package controllers

import (
	"net/http"
	"tamago/internal/api"
	"tamago/internal/models"

	"github.com/stretchr/testify/assert"
)

func getUserModel(userID uint64) *models.User {
	var (
		email       = "test@test.com"
		firstName   = "First"
		lastName    = "Last"
		phoneNumber = "123-456-7890"
	)

	user := &models.User{Email: email, FirstName: firstName, LastName: lastName, PhoneNumber: phoneNumber}
	user.ID = userID

	return user
}

func getProductModel(productID uint64) *models.Product {
	name := "Service A"

	product := &models.Product{Name: name}
	product.ID = productID
	return product
}

func getUserProductModel(hostID, productID, userProductID uint64) models.UserProduct {
	var (
		title                 = "My Product"
		description           = "This is my product"
		active                = true
		totalPrice            = 42.0
		maxPrice              = 42.0
		minPrice              = 42.0
		maxMemberCount uint64 = 5
		minMemberCount uint64 = 1
	)

	userProduct := models.UserProduct{HostID: hostID, ProductID: productID,
		Title: title, Description: description, Active: active,
		TotalPrice: totalPrice, MaxPrice: maxPrice, MinPrice: minPrice, MaxMemberCount: maxMemberCount,
		MinMemberCount: minMemberCount}
	userProduct.ID = userProductID
	return userProduct
}

func (s *ControllerSuite) TestGetUserProductsHappyPath() {
	var (
		productID     uint64 = 1
		hostID        uint64 = 1
		userProductID uint64 = 1
	)

	hostModel := getUserModel(hostID)
	productModel := getProductModel(productID)
	userProductModel := getUserProductModel(hostID, productID, userProductID)
	userProductModels := []models.UserProduct{userProductModel}

	userProduct := &api.UserProduct{
		ID:          userProductID,
		Host:        &api.UserProfile{ID: hostID, FirstName: hostModel.FirstName},
		Product:     &api.Product{ID: productID, Name: productModel.Name},
		Title:       userProductModel.Title,
		Description: userProductModel.Description,
		Active:      userProductModel.Active,
		Price: &api.UserProduct_Price{Total: userProductModel.TotalPrice,
			Min: userProductModel.MinPrice, Max: userProductModel.MaxPrice},
		MemberCount: &api.UserProduct_MemberCount{Min: userProductModel.MinMemberCount,
			Max: userProductModel.MaxMemberCount},
	}
	userProducts := []*api.UserProduct{userProduct}

	expectedResponse := &api.GetUserProductsResponse{UserProducts: userProducts}

	s.DAO.On("GetUserProductsByUserID", hostID).Return(userProductModels, nil)
	s.DAO.On("GetProductByProductID", productID).Return(productModel, nil)
	s.DAO.On("GetUserByID", hostID).Return(hostModel, nil)

	data, code := s.GET("/api/my_products", true)

	assert.Equal(s.T(), http.StatusOK, code)
	response := &api.GetUserProductsResponse{}
	s.UnmarshalResponse(data, response)
	assert.EqualValues(s.T(), expectedResponse, response)
}
