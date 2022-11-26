package controllers

import (
	"errors"
	fmt "fmt"
	"net/http"
	"tamago/internal/api"
	"tamago/internal/models"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

const password = "password"

func getValidUser() *models.User {
	var (
		userID      uint64 = 1
		email              = "test@test.com"
		firstName          = "First"
		lastName           = "Last"
		phoneNumber        = "123-456-7890"
	)

	user := &models.User{Email: email, FirstName: firstName, LastName: lastName, PhoneNumber: phoneNumber}
	user.ID = userID

	return user
}

func (s *ControllerSuite) TestGetUserProfileByIDHappyPath() {
	// setup
	user := getValidUser()

	userProfile := &api.UserProfile{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
	}
	expectedResponse := &api.GetUserProfileByIDResponse{UserProfile: userProfile}

	s.DAO.On("GetUserByID", user.ID).Return(user, nil)

	// action
	data, code := s.GET(fmt.Sprintf("/api/users/%v/profile", user.ID), true)

	// assert
	assert.Equal(s.T(), http.StatusOK, code)
	response := &api.GetUserProfileByIDResponse{}
	s.UnmarshalResponse(data, response)
	assert.EqualValues(s.T(), expectedResponse, response)
}

func (s *ControllerSuite) TestGetUserProfileByIDWhenFailedToQueryUser() {
	// setup
	var userID uint64 = 1
	s.DAO.On("GetUserByID", userID).Return(&models.User{}, errors.New("DB error"))

	// action
	data, code := s.GET(fmt.Sprintf("/api/users/%v/profile", userID), true)

	// assert
	assert.Equal(s.T(), http.StatusInternalServerError, code)
	err := s.UnmarshalErrorResponse(data)
	assert.Equal(s.T(), UnknownServerError, err.Code)
}

func (s *ControllerSuite) TestGetUserProfileByIDPolicyWhenUserIsNotLoggedIn() {
	data, code := s.GET(fmt.Sprintf("/api/users/%v/profile", uint64(1)), false)
	assert.Equal(s.T(), http.StatusUnauthorized, code)
	err := s.UnmarshalErrorResponse(data)
	assert.Equal(s.T(), LoginRequiredError, err.Code)
}

func (s *ControllerSuite) TestGetUserProfileByIDPolicyWhenUserIdNotMatch() {
	data, code := s.GET(fmt.Sprintf("/api/users/%v/profile", uint64(2)), true)
	assert.Equal(s.T(), http.StatusUnauthorized, code)
	err := s.UnmarshalErrorResponse(data)
	assert.Equal(s.T(), UnauthorizedError, err.Code)
}

func (s *ControllerSuite) TestSignUpUserHappyPath() {
	// setup
	user := getValidUser()

	request := &api.SignUpUserRequest{Email: user.Email, Password: password}

	s.DAO.On("GetUserByEmail", user.Email).Return(&models.User{}, gorm.ErrRecordNotFound)
	s.DAO.On("SignUpUser", user.Email, password).Return(user, nil)

	// action
	data, code := s.POST("/api/users/signup", request, false)

	// assert
	assert.Equal(s.T(), http.StatusOK, code)
	response := &api.SignUpUserResponse{}
	s.UnmarshalResponse(data, response)
	assert.Equal(s.T(), true, response.Success)
}

func (s *ControllerSuite) TestSignUpUserWhenEmailExists() {
	// setup
	user := getValidUser()

	request := &api.SignUpUserRequest{Email: user.Email, Password: password}

	s.DAO.On("GetUserByEmail", user.Email).Return(user, nil)

	// action
	data, code := s.POST("/api/users/signup", request, false)

	// assert
	assert.Equal(s.T(), http.StatusBadRequest, code)
	err := s.UnmarshalErrorResponse(data)
	assert.Equal(s.T(), FailedSignUpError, err.Code)
}

func (s *ControllerSuite) TestSignUpUserWhenFailedToCreateUser() {
	// setup
	email := "test@test.com"

	request := &api.SignUpUserRequest{Email: email, Password: password}

	s.DAO.On("GetUserByEmail", email).Return(&models.User{}, gorm.ErrRecordNotFound)
	s.DAO.On("SignUpUser", email, password).Return(&models.User{}, errors.New("DB error"))

	// action
	data, code := s.POST("/api/users/signup", request, false)

	// assert
	assert.Equal(s.T(), http.StatusInternalServerError, code)
	err := s.UnmarshalErrorResponse(data)
	assert.Equal(s.T(), UnknownServerError, err.Code)
}
