package controllers

import (
	"tamago/internal/api"
	"tamago/internal/context"
)

//GetProductsByBillingType returns products based on billing type in request, return a slice of product obj
func GetProductsByBillingType(r *context.RequestContext) (interface{}, error) {
	request, _ := r.GetRequest().(*api.GetProductsByBillingTypeRequest)
	billingType := request.GetBillingType()

	products, err := r.GetDAO().GetProductsByBillingType(billingType)

	if err != nil {
		return api.GetProductByIdResponse{}, err
	}

	var rstProducts []*api.Product
	for _, p := range *products {
		np := &api.Product{
			ID:             p.ID,
			Name:           p.Name,
			BillingType:    api.Product_BillingType(p.BillingType),
			IsFixedPrice:   p.IsFixedPrice,
			MaxMemberCount: p.MaxMemberCount,
		}
		rstProducts = append(rstProducts, np)
	}

	return api.GetProductsByBillingTypeResponse{Products: rstProducts}, nil
}
