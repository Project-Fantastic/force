package controllers

import (
	"tamago/internal/api"
	"tamago/internal/context"
)

// GetUserProducts returns a list of a user's joined or hosted products
func GetUserProducts(r *context.RequestContext) (interface{}, error) {
	userID := r.GetUserID()
	userProducts, err := r.GetDAO().GetUserProductsByUserID(userID)

	if err != nil {
		return api.GetUserProductsResponse{}, err
	}

	var rstUserProducts []*api.UserProduct

	for _, up := range userProducts {
		hostID := up.HostID
		productID := up.ProductID

		host, err := r.GetDAO().GetUserByID(hostID)

		if err != nil {
			continue
		}

		product, err := r.GetDAO().GetProductByProductID(productID)
		if err != nil {
			continue
		}

		userProduct := &api.UserProduct{
			ID:          up.ID,
			Host:        &api.UserProfile{ID: hostID, FirstName: host.FirstName},
			Product:     &api.Product{ID: productID, Name: product.Name},
			Title:       up.Title,
			Description: up.Description,
			Active:      up.Active,
			Price:       &api.UserProduct_Price{Total: up.TotalPrice, Min: up.MinPrice, Max: up.MaxPrice},
			MemberCount: &api.UserProduct_MemberCount{Min: up.MinMemberCount, Max: up.MaxMemberCount},
		}

		rstUserProducts = append(rstUserProducts, userProduct)
	}

	return api.GetUserProductsResponse{UserProducts: rstUserProducts}, nil
}
