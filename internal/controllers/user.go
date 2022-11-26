package controllers

import (
	"log"
	"net/http"
	"tamago/internal/api"
	"tamago/internal/context"

	"github.com/jinzhu/gorm"
)

// GetUserProfileByIDPolicy sets policy check for GetUserProfileByID endpoint.
func GetUserProfileByIDPolicy(r *context.RequestContext) (interface{}, error) {
	currentUserID := r.GetUserID()
	request, _ := r.GetRequest().(*api.GetUserProfileByIDRequest)
	// check if the logged in account has access to view the user profile
	if currentUserID != request.UserID {
		// TODO: allow to access other users' profile
		return nil, Errors[UnauthorizedError]
	}

	return nil, nil
}

// GetUserProfileByID returns a user's profile by the given ID.
func GetUserProfileByID(r *context.RequestContext) (interface{}, error) {
	request, _ := r.GetRequest().(*api.GetUserProfileByIDRequest)
	user, err := r.GetDAO().GetUserByID(request.UserID)
	if err != nil {
		return api.GetUserProfileByIDResponse{}, err
	}
	// TODO: hide email/phone_number when necessary
	userProfile := &api.UserProfile{ID: request.UserID, Email: user.Email, FirstName: user.FirstName,
		LastName: user.LastName, PhoneNumber: user.PhoneNumber}
	return api.GetUserProfileByIDResponse{UserProfile: userProfile}, nil
}

// SignUpUser signs up a new user.
func SignUpUser(r *context.RequestContext) (interface{}, error) {
	request, _ := r.GetRequest().(*api.SignUpUserRequest)
	dao := r.GetDAO()

	if _, err := dao.GetUserByEmail(request.Email); !gorm.IsRecordNotFoundError(err) {
		return api.SignUpUserResponse{Success: false},
			RequestError{"Failed to sign up", FailedSignUpError, http.StatusBadRequest}
	}
	user, err := dao.SignUpUser(request.Email, request.Password)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return api.SignUpUserResponse{Success: false}, err
	}
	// TODO: send an email for verification
	log.Println(user.ID)
	return api.SignUpUserResponse{Success: true}, nil
}
