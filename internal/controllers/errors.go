package controllers

import "net/http"

const (
	// authentication authorization errors
	LoginRequiredError = 10000 // Accessing resource requires logged in account
	UnauthorizedError  = 10001 // Accessing resource requires some level of permission
	// request validation errors
	RequestValidationError = 20000
	MissingRequiredError   = 20001
	InvalidEmailError      = 20002
	// denied access to some resources errors
	NotExistError = 30000 // Requested resource does not exist
	// failed requests due to user errors
	FailedSignUpError = 40000
	// server side errors
	UnknownServerError  = 50000 // Unknown server error
	FailedRequestError  = 50001
	FailedResponseError = 50002
)

var (
	loginRequiredError = RequestError{"Login is required", LoginRequiredError, http.StatusUnauthorized}
	unauthorizedError  = RequestError{"Don't have permissions to access this resource", UnauthorizedError,
		http.StatusUnauthorized}
	notExistError      = RequestError{"Requested resource does not exist", NotExistError, http.StatusNotFound}
	unknownServerError = RequestError{"Server currently is experiencing issues", UnknownServerError,
		http.StatusInternalServerError}
	failedRequestError = RequestError{"Failed to handle the request", FailedRequestError,
		http.StatusInternalServerError}
	failedResponseError = RequestError{"Failed to response the request", FailedResponseError,
		http.StatusInternalServerError}
)

// Errors is a map of common errors
var Errors = map[int]RequestError{
	LoginRequiredError:  loginRequiredError,
	UnauthorizedError:   unauthorizedError,
	NotExistError:       notExistError,
	UnknownServerError:  unknownServerError,
	FailedRequestError:  failedRequestError,
	FailedResponseError: failedResponseError,
}

type RequestError struct {
	Message    string `json:"message"`
	Code       int    `json:"code"`
	HTTPStatus int    `json:"-"`
}

func (r RequestError) Error() string {
	return r.Message
}
