package dao

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DAOSuite struct {
	suite.Suite
	DB   *gorm.DB
	Mock sqlmock.Sqlmock
	DAO  *DAO
}

func (s *DAOSuite) SetupTest() {
	mockedDB, mock, err := sqlmock.New()
	require.NoError(s.T(), err)

	db, err := gorm.Open("postgres", mockedDB)
	require.NoError(s.T(), err)

	db.LogMode(true)

	s.DB = db
	s.Mock = mock
	s.DAO = NewDAO(s.DB)
}

func (s *DAOSuite) TearDownTest() {
	s.DB.Close()
}

func (s *DAOSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.Mock.ExpectationsWereMet())
}

func TestDAOSuite(t *testing.T) {
	suite.Run(t, new(DAOSuite))
}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}
