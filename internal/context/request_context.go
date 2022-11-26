package context

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
	"strings"
	"tamago/internal/dao"
	"tamago/internal/sessions"

	"github.com/gin-gonic/gin"
	gsessions "github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
)

// Resource is an interface of app resources
type Resource interface {
	GetDatabase() *gorm.DB
	GetSession() sessions.SessionIface
	GetDAO() dao.DataAccessIface
}

// RequestContext is a container of gin.Context, DB connection, Session etc.
type RequestContext struct {
	Context      *gin.Context
	PolicyOutput interface{}
	request      reflect.Value
	resource     Resource
}

// NewRequestContext creates a new RequestContext
func NewRequestContext(context *gin.Context, resource Resource, request reflect.Value) *RequestContext {
	return &RequestContext{Context: context, resource: resource, request: request}
}

func (r *RequestContext) ConvertRequest() error {
	request := r.request.Elem()

	if request.Kind() != reflect.Struct {
		return errors.New("request is not a proper struct")
	}

	params := r.Context.Params
	for _, param := range params {
		if err := setReflectedStruct(request, param.Key, param.Value); err != nil {
			return err
		}
	}
	queryParams := r.Context.Request.URL.Query()
	for key, values := range queryParams {
		if err := setReflectedStruct(request, key, strings.Join(values[:], ",")); err != nil {
			return err
		}
	}
	if r.Context.ContentType() == gin.MIMEJSON {
		if err := r.decodeJSON(); err != nil {
			return err
		}
	}
	return nil
}

// GetRequest returns ProtoBuf defined API request object
func (r *RequestContext) GetRequest() interface{} {
	return r.request.Interface()
}

// GetDatabase returns a DB object
func (r *RequestContext) GetDatabase() *gorm.DB {
	return r.resource.GetDatabase()
}

func (r *RequestContext) GetDAO() dao.DataAccessIface {
	return r.resource.GetDAO()
}

// GetSession returns a raw sessions.Session object
func (r *RequestContext) GetSession() *gsessions.Session {
	session, err := r.resource.GetSession().Get(r.Context.Request, "session")
	if err != nil {
		log.Println(err)
	}
	return session
}

// IsUserLoggedIn returns if the current session is logged in
func (r *RequestContext) IsUserLoggedIn() bool {
	session := r.GetSession()
	_, ok := session.Values["user_id"]
	return ok
}

// GetUserID returns the current logged in user ID
func (r *RequestContext) GetUserID() uint64 {
	session := r.GetSession()
	userID, ok := session.Values["user_id"]
	if !ok {
		return 0
	}
	return userID.(uint64)
}

// SaveUserID stores userID into the session
func (r *RequestContext) SaveUserID(userID uint64) error {
	session := r.GetSession()
	session.Values["user_id"] = userID
	return r.resource.GetSession().Save(r.Context.Request, r.Context.Writer, session)
}

// Logout clears the current user from the session
func (r *RequestContext) Logout() error {
	session := r.GetSession()
	session.Options.MaxAge = -1
	return r.resource.GetSession().Save(r.Context.Request, r.Context.Writer, session)
}

func (r *RequestContext) decodeJSON() error {
	decoder := json.NewDecoder(r.Context.Request.Body)

	if err := decoder.Decode(r.GetRequest()); err != nil && err != io.EOF {
		return err
	}
	return nil
}

func setReflectedStruct(r reflect.Value, key string, value string) error {
	field := r.FieldByName(strings.Join(convert(strings.Split(key, "_")), ""))
	if !field.IsValid() || !field.CanSet() {
		log.Printf("not able to set %s \n", key)
		return nil
	}

	return setReflectedField(field, value)
}

func setReflectedField(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.Int32:
		return setReflectedInt(field, value, 32)
	case reflect.Int64:
		return setReflectedInt(field, value, 64)
	case reflect.Uint32:
		return setReflectedUint(field, value, 32)
	case reflect.Uint64:
		return setReflectedUint(field, value, 64)
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		return setReflectedBool(field, value)
	case reflect.Float32:
		return setReflectedFloat(field, value, 32)
	case reflect.Float64:
		return setReflectedFloat(field, value, 64)
	case reflect.Slice:
		return setReflectedSlice(field, strings.Split(value, ","))
	default:
		return fmt.Errorf("field type: %v is not supported", field.Kind())
	}
	return nil
}

func setReflectedSlice(field reflect.Value, values []string) error {
	slice := reflect.MakeSlice(field.Type(), len(values), len(values))
	for i, value := range values {
		if err := setReflectedField(slice.Index(i), value); err != nil {
			return err
		}
	}
	field.Set(slice)
	return nil
}

func setReflectedInt(field reflect.Value, value string, bitSize int) error {
	if value == "" {
		value = "0"
	}
	v, err := strconv.ParseInt(value, 10, bitSize)
	if err == nil {
		field.SetInt(v)
	}
	return err
}

func setReflectedUint(field reflect.Value, value string, bitSize int) error {
	if value == "" {
		value = "0"
	}
	v, err := strconv.ParseUint(value, 10, bitSize)
	if err == nil {
		field.SetUint(v)
	}
	return err
}

func setReflectedBool(field reflect.Value, value string) error {
	if value == "" {
		value = "false"
	}
	v, err := strconv.ParseBool(value)
	if err == nil {
		field.SetBool(v)
	}
	return err
}

func setReflectedFloat(field reflect.Value, value string, bitSize int) error {
	if value == "" {
		value = "0.0"
	}
	v, err := strconv.ParseFloat(value, bitSize)
	if err == nil {
		field.SetFloat(v)
	}
	return err
}

func convert(strs []string) []string {
	tempStrs := make([]string, len(strs))
	for i, v := range strs {
		if strings.ToLower(v) == "id" {
			tempStrs[i] = strings.ToUpper(v)
		} else {
			tempStrs[i] = strings.Title(v)
		}
	}
	return tempStrs
}
