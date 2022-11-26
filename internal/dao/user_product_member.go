package dao

import "tamago/internal/models"

//GetUserProductMembersByUserProductID returns a slice of members for a specific user product
func (d *DAO) GetUserProductMembersByUserProductID(upID uint64) (*[]models.UserProductMember, error) {
	member := &models.UserProductMember{}
	return member.GetUserProductMembersByUserProductID(d.db, upID)
}

//GetUserProductMemberByID returns a user product member by ID
func (d *DAO) GetUserProductMemberByID(id uint64) (*models.UserProductMember, error) {
	member := &models.UserProductMember{}
	return member.GetUserProductMemberByID(d.db, id)
}

//CreateUserProductMember creates a new user product member and returns it
func (d *DAO) CreateUserProductMember(
	userProductID uint64,
	userID uint64,
	isHost bool) (*models.UserProductMember, error) {
	m := &models.UserProductMember{}
	m.UserID = userID
	m.UserProductID = userProductID
	m.IsHost = isHost
	return m.CreateUserProductMember(d.db)
}
