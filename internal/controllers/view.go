package controllers

import (
	"reflect"
	"tamago/internal/context"
	"tamago/internal/dao"
	"tamago/internal/sessions"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type ViewHandler func(rc *context.RequestContext)

type View struct {
	engine  *gin.Engine
	db      *gorm.DB
	session sessions.SessionIface
	dao     dao.DataAccessIface
}

func NewView(engine *gin.Engine, db *gorm.DB, sessionStore sessions.SessionIface, dao dao.DataAccessIface) *View {
	return &View{engine: engine, db: db, session: sessionStore, dao: dao}
}

func (v *View) GetDatabase() *gorm.DB {
	return v.db
}

func (v *View) GetSession() sessions.SessionIface {
	return v.session
}

func (v *View) GetDAO() dao.DataAccessIface {
	return v.dao
}

func (v *View) Handle(method, path string, f ViewHandler) {
	v.engine.Handle(method, path, func(c *gin.Context) {
		requestContext := context.NewRequestContext(c, v, reflect.Zero(reflect.TypeOf(0)))
		f(requestContext)
	})
}

func (v *View) GET(path string, f ViewHandler) {
	v.Handle("GET", path, f)
}

func (v *View) POST(path string, f ViewHandler) {
	v.Handle("POST", path, f)
}
