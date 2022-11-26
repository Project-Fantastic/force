package dao

import (
	"log"

	"tamago/internal/models"
)

//GetUserProductsByProductID returns a slice of user products based on product id
func (d *DAO) GetUserProductsByProductID(productID uint64) (*[]models.UserProduct, error) {
	up := &models.UserProduct{}
	userProducts, err := up.GetUserProductsByProductID(d.db, productID)

	if err != nil {
		return nil, err
	}

	for _, up := range *userProducts {
		userProduct := &up //Because `up` is value, get memeory addr for value
		upID := userProduct.ID
		m := &models.UserProductMember{}
		members, err := m.GetUserProductMembersByUserProductID(d.db, upID)
		if err != nil {
			log.Printf("No member was found for user product: %d, with error: %s", upID, err)
			continue
		}
		userProduct.UserProductMembers = *members
	}

	return userProducts, err
}

//GetUserProductByUserProductID returns a user product by user product id
func (d *DAO) GetUserProductByUserProductID(id uint64) (*models.UserProduct, error) {
	up := &models.UserProduct{}
	up, err := up.GetUserProductByUserProductID(d.db, id)
	if err != nil {
		return nil, err
	}

	upm := &models.UserProductMember{}
	members, err := upm.GetUserProductMembersByUserProductID(d.db, id)
	if err != nil {
		return nil, err
	}
	up.UserProductMembers = *members

	br := &models.BillingRequest{}
	billingRequests, err := br.GetBillingRequestByUserProductID(d.db, id)
	if err != nil {
		return nil, err
	}
	up.BillingRequests = *billingRequests

	return up, nil
}

// GetUserProductsByUserID returns
func (d *DAO) GetUserProductsByUserID(userID uint64) ([]models.UserProduct, error) {
	u := &models.UserProduct{}
	userProducts, err := u.GetUserProductsByUserID(d.db, userID)
	return userProducts, err
}
