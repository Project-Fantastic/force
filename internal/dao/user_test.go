package dao

import (
	"tamago/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *DAOSuite) TestGetUserByIDHappyPath() {
	var (
		userID uint64 = 1
		email         = "test@test.com"
	)

	rows := sqlmock.NewRows([]string{"id", "email"}).AddRow(userID, email)

	s.Mock.ExpectQuery("SELECT").WillReturnRows(rows)

	user, err := s.DAO.GetUserByID(userID)

	expectedUser := &models.User{}
	expectedUser.ID = userID
	expectedUser.Email = email

	require.Nil(s.T(), err)
	assert.EqualValues(s.T(), user, expectedUser)
}

func (s *DAOSuite) TestGetUserByIDWhenUserNotExist() {
	var userID uint64 = 999

	s.Mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{}))

	user, err := s.DAO.GetUserByID(userID)

	require.Error(s.T(), err)
	assert.EqualValues(s.T(), user, &models.User{})
}

func (s *DAOSuite) TestGetUserByEmailHappyPath() {
	var (
		userID uint64 = 1
		email         = "test@test.com"
	)

	rows := sqlmock.NewRows([]string{"id", "email"}).AddRow(userID, email)

	s.Mock.ExpectQuery("SELECT").WillReturnRows(rows)

	user, err := s.DAO.GetUserByEmail(email)

	expectedUser := &models.User{}
	expectedUser.ID = userID
	expectedUser.Email = email

	require.Nil(s.T(), err)
	assert.EqualValues(s.T(), user, expectedUser)
}

func (s *DAOSuite) TestGetUserByEmailWhenEmailNotExist() {
	email := "non_existing@test.com"

	s.Mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{}))

	user, err := s.DAO.GetUserByEmail(email)

	require.Error(s.T(), err)
	assert.EqualValues(s.T(), user, &models.User{})
}

func (s *DAOSuite) TestCreateUserHappyPath() {
	var (
		email           = "test@test.com"
		password        = "password"
		userID   uint64 = 1
	)

	s.Mock.ExpectBegin()
	s.Mock.ExpectQuery("INSERT INTO \"users\"").
		WithArgs(AnyTime{}, AnyTime{}, nil, email, "", sqlmock.AnyArg(), "", "", "").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))
	s.Mock.ExpectQuery("SELECT").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow(0))
	s.Mock.ExpectCommit()

	user, err := s.DAO.SignUpUser(email, password)

	require.Nil(s.T(), err)
	assert.Equal(s.T(), userID, user.ID)
	assert.Equal(s.T(), email, user.Email)
	assert.NotEmpty(s.T(), password, user.Password)
	assert.NotEqual(s.T(), password, user.Password)
	assert.Empty(s.T(), user.PhoneNumber)
	assert.NotNil(s.T(), user.CreatedAt)
	assert.NotNil(s.T(), user.UpdatedAt)
}

func (s *DAOSuite) TestCreateUserOnValidationError() {
	email := "test@test.com"

	_, err := s.DAO.SignUpUser(email, "")

	require.Error(s.T(), err)
}

func (s *DAOSuite) TestCreateUserOnUniqueKeyError() {
	var (
		email    = "test@test.com"
		password = "password"
	)

	s.Mock.ExpectBegin()
	s.Mock.ExpectQuery("INSERT INTO \"users\"").
		WithArgs(AnyTime{}, AnyTime{}, nil, email, "", sqlmock.AnyArg(), "", "", "").
		WillReturnError(pq.Error{Message: "unique_violation"})
	s.Mock.ExpectRollback()

	_, err := s.DAO.SignUpUser(email, password)

	require.Error(s.T(), err)
}

func (s *DAOSuite) TestVerifyLoginHappyPath() {
	var (
		expectedUserID uint64 = 1
		email                 = "test@test.com"
		password              = "password"
	)

	hashedPassword, _ := models.HashPassword(password)

	rows := sqlmock.NewRows([]string{"id", "email", "password"}).AddRow(expectedUserID, email, hashedPassword)

	s.Mock.ExpectQuery("SELECT").WillReturnRows(rows)

	verified, userID := s.DAO.VerifyLogin(email, password)

	require.Equal(s.T(), true, verified)
	require.Equal(s.T(), expectedUserID, userID)
}

func (s *DAOSuite) TestVerifyLoginWhenEmailNotExist() {
	var (
		email    = "test@test.com"
		password = "password"
	)

	s.Mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{}))

	verified, userID := s.DAO.VerifyLogin(email, password)

	require.Equal(s.T(), false, verified)
	require.Equal(s.T(), uint64(0), userID)
}

func (s *DAOSuite) TestVerifyLoginWithMismatchedPassword() {
	var (
		expectedUserID uint64 = 1
		email                 = "test@test.com"
		password              = "password"
	)

	hashedPassword, _ := models.HashPassword(password)

	rows := sqlmock.NewRows([]string{"id", "email", "password"}).AddRow(expectedUserID, email, hashedPassword)

	s.Mock.ExpectQuery("SELECT").WillReturnRows(rows)

	verified, userID := s.DAO.VerifyLogin(email, "wrong_password")

	require.Equal(s.T(), false, verified)
	require.Equal(s.T(), uint64(0), userID)
}
