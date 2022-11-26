package dao

import "tamago/internal/models"

// GetUserByID returns a user by ID.
func (d *DAO) GetUserByID(userID uint64) (*models.User, error) {
	user := &models.User{}
	return user.GetUserByID(d.db, userID)
}

// GetUserByEmail returns a user by email.
func (d *DAO) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	return user.GetUserByEmail(d.db, email)
}

// SignUpUser signs up a new user with email and password.
// TODO: handle phone number, first and last name if needed.
func (d *DAO) SignUpUser(email string, password string) (*models.User, error) {
	user := &models.User{Email: email, Password: password}
	return user.CreateUser(d.db)
}

// VerifyLogin takes email and password and verify if the credential is correct
func (d *DAO) VerifyLogin(email, password string) (bool, uint64) {
	user := &models.User{}
	user, err := user.GetUserByEmail(d.db, email)
	if err != nil {
		return false, 0
	}
	if models.VerifyPassword(user.Password, password) {
		return true, user.ID
	}
	return false, 0
}
