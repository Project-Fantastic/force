package controllers

import (
	fmt "fmt"
	"log"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"tamago/internal/context"
	"tamago/internal/dao"
	"tamago/internal/sessions"

	"github.com/gin-gonic/gin"
	proto "github.com/gogo/protobuf/proto"
	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"
)

type RequestHandler func(r *context.RequestContext) (interface{}, error)

type PolicyHandler func(r *context.RequestContext) (interface{}, error)

type ValidatorHandler func(r *context.RequestContext) error

type API struct {
	engine    *gin.Engine
	db        *gorm.DB
	session   sessions.SessionIface
	dao       dao.DataAccessIface
	endpoints map[string]*Endpoint
}

type Endpoint struct {
	method        string
	path          string
	policy        PolicyHandler
	validators    []ValidatorHandler
	loginRequired bool
}

type RequestValidator interface {
	Validate(interface{}) error
}

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

const prefix = "/api/"

var requestValidator RequestValidator = &defaultValidator{}

func NewAPI(engine *gin.Engine, db *gorm.DB, sessionStore sessions.SessionIface, dao dao.DataAccessIface) *API {
	endpoints := make(map[string]*Endpoint)
	return &API{engine: engine, db: db, session: sessionStore, dao: dao, endpoints: endpoints}
}

func (a *API) GetDatabase() *gorm.DB {
	return a.db
}

func (a *API) GetSession() sessions.SessionIface {
	return a.session
}

func (a *API) GetDAO() dao.DataAccessIface {
	return a.dao
}

func (a *API) Handle(method, path string, f RequestHandler) *Endpoint {
	endpoint := &Endpoint{method: method, path: path, loginRequired: true}
	endpointKey := fmt.Sprintf("%s_%s", method, path)
	_, ok := a.endpoints[endpointKey]
	if ok {
		log.Printf("Endpoint %s exists\n", endpointKey)
		return endpoint
	}
	endpoint.validators = make([]ValidatorHandler, 0)
	a.endpoints[endpointKey] = endpoint
	a.engine.Handle(method, prefix+path, func(c *gin.Context) {
		handlerName := getRequestHandlerName(f)
		requestType := proto.MessageType(handlerName + "Request")
		responseType := proto.MessageType(handlerName + "Response")
		if requestType == nil || responseType == nil {
			log.Println("Request or response type is nil")
			handleError(c, Errors[FailedRequestError])
			return
		}
		request := reflect.New(requestType.Elem())
		requestContext := context.NewRequestContext(c, a, request)
		if err := requestContext.ConvertRequest(); err != nil {
			log.Printf("Failed to convert a request: %s\n", err)
			handleError(c, Errors[FailedRequestError])
			return
		}

		isUserLoggedIn := requestContext.IsUserLoggedIn()

		if endpoint.loginRequired && !isUserLoggedIn {
			handleError(c, Errors[LoginRequiredError])
			return
		}

		if endpoint.policy != nil {
			policyOutput, err := endpoint.policy(requestContext)
			if err != nil {
				handleError(c, err)
				return
			}
			requestContext.PolicyOutput = policyOutput
		}
		for _, validator := range endpoint.validators {
			err := validator(requestContext)
			if err != nil {
				handleError(c, err)
				return
			}
		}
		response, err := f(requestContext)
		if reflect.TypeOf(response) != responseType.Elem() {
			err = Errors[FailedResponseError]
		}
		if err != nil {
			handleError(c, err)
		} else {
			c.JSON(http.StatusOK, response)
		}
	})
	return endpoint
}

func (a *API) GET(path string, f RequestHandler) *Endpoint {
	return a.Handle("GET", path, f)
}

func (a *API) POST(path string, f RequestHandler) *Endpoint {
	return a.Handle("POST", path, f)
}

func (a *API) PUT(path string, f RequestHandler) *Endpoint {
	return a.Handle("PUT", path, f)
}

func (a *API) DELETE(path string, f RequestHandler) *Endpoint {
	return a.Handle("DELETE", path, f)
}

func (e *Endpoint) LoginRequired(isRequired bool) *Endpoint {
	e.loginRequired = isRequired
	return e
}

func (e *Endpoint) SetPolicy(f PolicyHandler) *Endpoint {
	e.policy = f
	return e
}

func (e *Endpoint) AddValidator(f ValidatorHandler) *Endpoint {
	e.validators = append(e.validators, f)
	return e
}

func (e *Endpoint) EnableDefaultValidation() *Endpoint {
	return e.AddValidator(defaultRequestValidator)
}

func defaultRequestValidator(r *context.RequestContext) error {
	return requestValidator.Validate(r.GetRequest())
}

func getRequestHandlerName(f RequestHandler) string {
	return "api." + trimFuncName(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
}

func (v *defaultValidator) Validate(request interface{}) error {
	v.once.Do(func() {
		v.validate = validator.New()
	})
	if err := v.validate.Struct(request); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				return RequestError{fmt.Sprintf("Field: %s is required", err.Field()),
					MissingRequiredError, http.StatusBadRequest}
			case "email":
				return RequestError{fmt.Sprintf("Field: %s is not a valid email address", err.Field()),
					InvalidEmailError, http.StatusBadRequest}
			default:
				return RequestError{fmt.Sprintf("Field: %s is not valid", err.Field()),
					RequestValidationError, http.StatusBadRequest}
			}
		}
	}
	return nil
}

func trimFuncName(fullName string) string {
	return strings.TrimPrefix(filepath.Ext(fullName), ".")
}

func handleError(c *gin.Context, e error) {
	err, ok := e.(RequestError)
	if ok {
		c.JSON(err.HTTPStatus, err)
		return
	} else if gorm.IsRecordNotFoundError(e) {
		err = Errors[NotExistError]
	} else {
		err = Errors[UnknownServerError]
	}
	c.JSON(err.HTTPStatus, err)
}
