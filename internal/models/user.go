package models

import (
	"errors"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User is a model of users.
type User struct {
	BaseModel
	Email                  string        `gorm:"type:varchar(255);not null;unique_index"`
	PhoneNumber            string        `gorm:"type:varchar(20)"`
	Password               string        `gorm:"type:varchar(512);not null"`
	FirstName              string        `gorm:"type:varchar(100)"`
	LastName               string        `gorm:"type:varchar(100)"`
	ProfilePicture         string        `gorm:"type:varchar(1024)"`
	Status                 uint          `gorm:"not null;default:'0'"`
	UserProducts           []UserProduct `gorm:"foreignkey:HostID"`
	UserProductMembers     []UserProductMember
	UserThirdPartyPayments []UserThirdPartyPayment
	UserProductAccounts    []UserProductAccount
	Notifications          []Notification
}

// HashPassword is used to generate hashed password with salt.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword compares a hashed password with a plain text to check if they match.
func VerifyPassword(hashedPassword, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}
	return true
}

func (u *User) prepare() {
	currentTime := time.Now().UTC()
	u.ID = 0
	u.Email = strings.TrimSpace(u.Email)
	u.PhoneNumber = strings.TrimSpace(u.PhoneNumber)
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.CreatedAt = currentTime
	u.UpdatedAt = currentTime
	u.Status = 0
}

func (u *User) validate() error {
	if u.Email == "" {
		return errors.New("Email is required")
	}

	if u.Password == "" {
		return errors.New("Password is required")
	}

	return nil
}

func (u *User) beforeSave() error {
	password, err := HashPassword(u.Password)
	if err != nil {
		return err
	}

	u.Password = password
	return nil
}

// GetUserByID returns a User model by a given ID.
func (u *User) GetUserByID(db *gorm.DB, userID uint64) (*User, error) {
	err := db.First(u, userID).Error
	return u, err
}

// GetUserByEmail returns a User model by a given Email.
func (u *User) GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	err := db.Where("email = ?", email).First(u).Error
	return u, err
}

// CreateUser create a new User.
func (u *User) CreateUser(db *gorm.DB) (*User, error) {
	u.prepare()
	err := u.validate()
	if err != nil {
		return &User{}, err
	}
	err = u.beforeSave()
	if err != nil {
		return &User{}, err
	}
	err = db.Create(u).Error
	return u, err
}
