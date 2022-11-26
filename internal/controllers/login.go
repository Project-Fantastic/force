package controllers

import (
	fmt "fmt"
	"log"
	"net/http"
	"tamago/internal/context"

	"github.com/gin-gonic/gin"
)

type LoginForm struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

func GetLoginView(r *context.RequestContext) {
	r.Context.HTML(http.StatusOK, "login.tmpl", gin.H{})
}

func PostLoginView(r *context.RequestContext) {
	var loginForm LoginForm
	redirectPath := "/login"

	if err := r.Context.ShouldBind(&loginForm); err == nil {
		verified, userID := r.GetDAO().VerifyLogin(loginForm.Email, loginForm.Password)
		if verified {
			err := r.SaveUserID(userID)
			if err == nil {
				redirectPath = "/"
			}
		}
	}
	r.Context.Redirect(http.StatusFound, redirectPath)
}

func GetLogoutView(r *context.RequestContext) {
	if err := r.Logout(); err != nil {
		log.Println(fmt.Sprintf("Failed to logout: %v", err))
	}
	r.Context.Redirect(http.StatusFound, "/login")
}
