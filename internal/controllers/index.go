package controllers

import (
	"tamago/internal/context"

	"github.com/gin-gonic/gin"
)

// GetIndexView renders Index page.
func GetIndexView(r *context.RequestContext) {
	r.Context.JSON(200, gin.H{"logged_in": r.IsUserLoggedIn(), "user_id": r.GetUserID()})
}
