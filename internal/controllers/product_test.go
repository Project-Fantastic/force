package controllers

import (
	fmt "fmt"
	"net/http"
	"tamago/internal/api"
	"tamago/internal/models"

	"github.com/stretchr/testify/assert"
)

func getValidProduct() *models.Product {
	var (
		productID      uint64 = 1
		name                  = "youtube"
		billingType           = api.Product_BILLING_TYPE_RECURRING
		isFixedPrice          = true
		maxMemberCount uint64 = 5
	)

	p := &models.Product{
		Name:           name,
		BillingType:    int32(billingType),
		IsFixedPrice:   isFixedPrice,
		MaxMemberCount: maxMemberCount,
	}
	p.ID = productID
	return p
}

func (s *ControllerSuite) TestGetProductsByBillingType() {
	// setup
	pModel := getValidProduct()

	productRestModel := &api.Product{
		ID:             pModel.ID,
		Name:           pModel.Name,
		BillingType:    api.Product_BillingType(pModel.BillingType),
		IsFixedPrice:   pModel.IsFixedPrice,
		MaxMemberCount: pModel.MaxMemberCount,
	}

	expectedResponse := &api.GetProductsByBillingTypeResponse{Products: []*api.Product{productRestModel}}

	s.DAO.On("GetProductsByBillingType", productRestModel.GetBillingType()).Return(&([]models.Product{*pModel}), nil)

	// action
	data, code := s.GET(fmt.Sprintf("/api/products?billing_type=%v", pModel.BillingType), false)

	// assert
	assert.Equal(s.T(), http.StatusOK, code)
	response := &api.GetProductsByBillingTypeResponse{}
	s.UnmarshalResponse(data, response)
	assert.EqualValues(s.T(), expectedResponse, response)
}
