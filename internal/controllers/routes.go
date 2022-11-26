package controllers

// CreateViewRoutes adds new routes for Views
func (v *View) CreateViewRoutes() {
	v.GET("/", GetIndexView)
	v.GET("/login", GetLoginView)
	v.POST("/login", PostLoginView)
	v.GET("/logout", GetLogoutView)
}

// CreateAPIRoutes adds new routes for API
func (a *API) CreateAPIRoutes() {
	// User APIs
	a.GET("/users/:user_id/profile", GetUserProfileByID).
		SetPolicy(GetUserProfileByIDPolicy)
	a.POST("/users/signup", SignUpUser).LoginRequired(false).EnableDefaultValidation()
	// Products APIs
	a.GET("/products", GetProductsByBillingType).LoginRequired(false).EnableDefaultValidation()
	// User Products APIs
	a.GET("/my_products", GetUserProducts).LoginRequired(true).EnableDefaultValidation()
}
