package sessions

import (
	"net/http"
	"time"

	gsessions "github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/wader/gormstore"
)

// 30 days
const age = 24 * 60 * 60 * 30

// 7 days
const hours = 168

// SessionIface is an interface for session ops.
type SessionIface interface {
	Get(*http.Request, string) (*gsessions.Session, error)
	Save(*http.Request, http.ResponseWriter, *gsessions.Session) error
}

// NewDBSession creates a global object that stores SQL-DB based sessions.
func NewDBSession(db *gorm.DB, secretKey string) *gormstore.Store {
	session := gormstore.NewOptions(
		db,
		gormstore.Options{TableName: "sessions", SkipCreateTable: true},
		[]byte(secretKey),
		[]byte(nil),
	)
	session.SessionOpts.MaxAge = age
	return session
}

// CleanDBSession cleans up expired sessions periodically.
func CleanDBSession(store *gormstore.Store, quit <-chan struct{}) {
	store.PeriodicCleanup(hours*time.Hour, quit)
}
