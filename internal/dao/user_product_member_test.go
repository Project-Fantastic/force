package dao

import (
	"tamago/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *DAOSuite) TestGetUserProductMembersByUserProductIDHappyPath() {
	var (
		memberID      uint64 = 1
		userProductID uint64 = 1
		isHost               = true
	)

	rows := sqlmock.NewRows([]string{"id", "user_product_id", "is_host"}).AddRow(memberID, userProductID, isHost)
	s.Mock.ExpectQuery("^SELECT").WillReturnRows(rows)

	members, err := s.DAO.GetUserProductMembersByUserProductID(userProductID)

	require.Nil(s.T(), err, nil)
	assert.Equal(s.T(), 1, len(*members))

	expectedMember := &models.UserProductMember{}
	expectedMember.ID = memberID
	expectedMember.UserProductID = userProductID
	expectedMember.IsHost = isHost

	assert.Equal(s.T(), expectedMember, &(*members)[0])
}

func (s *DAOSuite) TestGetUserProductMemberByIDHappyPath() {
	var (
		upmID         uint64 = 1
		userID        uint64 = 1
		userProductID uint64 = 1
		isHost               = true
	)

	rows := sqlmock.
		NewRows([]string{"id", "user_id", "user_product_id", "is_host"}).
		AddRow(upmID, userID, userProductID, isHost)
	s.Mock.ExpectQuery("SELECT").WillReturnRows(rows)

	members, err := s.DAO.GetUserProductMemberByID(upmID)

	require.Nil(s.T(), err, nil)
	expectedMember := &models.UserProductMember{}
	expectedMember.ID = upmID
	expectedMember.UserID = userID
	expectedMember.UserProductID = userProductID
	expectedMember.IsHost = isHost

	assert.Equal(s.T(), expectedMember, members)
}

func (s *DAOSuite) TestGetUserProductMemberByIDWhenNoMemberFound() {
	var upmID uint64 = 999

	s.Mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{}))

	upm, err := s.DAO.GetUserProductMemberByID(upmID)

	require.Error(s.T(), err)
	assert.Equal(s.T(), &models.UserProductMember{}, upm)
}

func (s *DAOSuite) TestCreateUserProductMemberHappyPath() {
	var (
		upmID         uint64 = 1
		userID        uint64 = 1
		userProductID uint64 = 1
		isHost               = true
	)

	s.Mock.ExpectBegin()
	s.Mock.ExpectQuery("INSERT INTO \"user_product_members\"").
		WithArgs(AnyTime{}, AnyTime{}, nil, userID, userProductID, isHost).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(upmID))
	s.Mock.ExpectQuery("SELECT").
		WithArgs(upmID).
		WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow(0))
	s.Mock.ExpectCommit()

	member, err := s.DAO.CreateUserProductMember(userProductID, userID, isHost)

	require.Nil(s.T(), err, nil)
	assert.Equal(s.T(), userProductID, member.UserProductID)
	assert.Equal(s.T(), userID, member.UserID)
	assert.Equal(s.T(), isHost, member.IsHost)
	assert.Empty(s.T(), member.Status)
	assert.NotNil(s.T(), member.CreatedAt)
	assert.NotNil(s.T(), member.UpdatedAt)
}
