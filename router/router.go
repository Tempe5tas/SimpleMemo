package router

import "github.com/gin-gonic/gin"

func Init() *gin.Engine {
	r := gin.Default()
	r.POST("/user/register", UserRegister)
	r.POST("/user/login", UserLogin)
	return r
}
