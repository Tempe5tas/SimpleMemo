package router

import "github.com/gin-gonic/gin"

func init() {
	r := gin.Default()

	r.Run(":8080")
}
