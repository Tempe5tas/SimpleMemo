package router

import (
	"SimpleMemo/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func MemoCreate(c *gin.Context) {
	var memo *model.Memo
	if err := c.ShouldBind(&memo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}
}
