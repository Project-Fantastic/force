package controllers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"tamago/internal/mocks"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/gogo/protobuf/jsonpb"
	proto "github.com/gogo/protobuf/proto"
	gsessions "github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const loggedInUserID uint64 = 1

type ControllerSuite struct {
	suite.Suite
	Engine  *gin.Engine
	DB      *gorm.DB
	Session *mocks.SessionIface
	API     *API
	DAO     *mocks.DataAccessIface
}

func (s *ControllerSuite) SetupTest() {
	s.Engine = gin.New()
	mockedDB, _, err := sqlmock.New()
	require.NoError(s.T(), err)
	db, err := gorm.Open("postgres", mockedDB)
	require.NoError(s.T(), err)
	s.DB = db
	s.DAO = new(mocks.DataAccessIface)
	s.Session = new(mocks.SessionIface)
	s.API = NewAPI(s.Engine, s.DB, s.Session, s.DAO)

	// register the routes
	s.API.CreateAPIRoutes()
}

func (s *ControllerSuite) TearDownTest() {
	s.DB.Close()
}

func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(ControllerSuite))
}

func (s *ControllerSuite) Request(method, path string, body io.Reader, isLoggedIn bool) (io.Reader, int) {
	req := httptest.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	store := gsessions.NewCookieStore([]byte("mock-cookie"))
	session := gsessions.NewSession(store, "session")
	if isLoggedIn {
		session.Values["user_id"] = loggedInUserID
	}
	err := session.Save(req, resp)
	require.Nil(s.T(), err)
	s.Session.On("Get", req, "session").Return(session, nil)

	s.Engine.ServeHTTP(resp, req)

	return resp.Body, resp.Code
}

func (s *ControllerSuite) RequestWithProto(method, path string, pb proto.Message, isLoggedIn bool) (io.Reader, int) {
	m := &jsonpb.Marshaler{OrigName: true}
	data, _ := m.MarshalToString(pb)
	return s.Request(method, path, strings.NewReader(data), isLoggedIn)
}

func (s *ControllerSuite) GET(path string, isLoggedIn bool) (io.Reader, int) {
	return s.Request("GET", path, nil, isLoggedIn)
}

func (s *ControllerSuite) POST(path string, pb proto.Message, isLoggedIn bool) (io.Reader, int) {
	return s.RequestWithProto("POST", path, pb, isLoggedIn)
}

func (s *ControllerSuite) UnmarshalResponse(body io.Reader, pb proto.Message) {
	err := jsonpb.Unmarshal(body, pb)
	require.Nil(s.T(), err)
}

func (s *ControllerSuite) UnmarshalErrorResponse(body io.Reader) *RequestError {
	data, err := ioutil.ReadAll(body)
	require.Nil(s.T(), err)
	requestError := &RequestError{}
	err = json.Unmarshal(data, requestError)
	require.Nil(s.T(), err)
	return requestError
}
