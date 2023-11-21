package router

import (
	"SimpleMemo/middleware"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	r := gin.Default()
	// User login and register routers
	loginGroup := r.Group("user")
	{
		loginGroup.POST("register", UserRegister)
		loginGroup.POST("login", UserLogin)
	}
	// User info and edit routers
	userGroup := r.Group("user")
	{
		userGroup.Use(middleware.ValidateToken)
		userGroup.GET("profile", UserProfile)
		userGroup.PUT("profile/update", UserUpdate)
	}
	// Memo CRUD routers
	memoGroup := r.Group("memo")
	{
		memoGroup.Use(middleware.ValidateToken)
		memoGroup.POST("create", MemoCreate)
		memoGroup.GET("list")
		memoGroup.PUT("update")
		memoGroup.DELETE("delete")
	}
	return r
}
