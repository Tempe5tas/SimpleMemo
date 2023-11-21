package router

import (
	"SimpleMemo/model"
	response "SimpleMemo/serializer"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func MemoCreate(c *gin.Context) {
	// Get object from request body
	var memoForm model.MemoCreate
	if err := c.ShouldBind(&memoForm); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code: http.StatusBadRequest,
			Msg:  "invalid memo format",
			Data: err.Error(),
		})
		return
	}
	// Verify user ID
	ID, ok := c.Get("ID")
	var existed int64
	model.DB.Where("ID = ?", ID).First(&model.User{}).Count(&existed)
	if existed != 1 || !ok {
		c.JSON(http.StatusBadRequest, response.Response{
			Code: http.StatusBadRequest,
			Msg:  "no such user",
		})
		return
	}
	// Validate time format
	memoTime, err := time.Parse("2006-01-02 15:04:05", memoForm.Time)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code: http.StatusBadRequest,
			Msg:  "invalid time format",
			Data: err.Error(),
		})
		return
	}
	// Create memo
	var user model.User
	model.DB.Take(&user, ID)
	memo := &model.Memo{
		User:    user,
		Title:   memoForm.Title,
		Time:    memoTime,
		Status:  memoForm.Status,
		Content: memoForm.Content,
	}
	model.DB.Create(&memo)
	c.JSON(http.StatusOK, response.Response{
		Code: http.StatusOK,
		Msg:  "memo creation successful",
		Data: map[string]any{"ID": memo.ID},
	})
	return
}

func MemoRetrieve(c *gin.Context) {
	// Verify user ID
	ID, ok := c.Get("ID")
	var existed int64
	model.DB.Where("ID = ?", ID).First(&model.User{}).Count(&existed)
	if existed != 1 || !ok {
		c.JSON(http.StatusBadRequest, response.Response{
			Code: http.StatusBadRequest,
			Msg:  "no such user",
		})
		return
	}
	// Retrieve memos from certain user
	var user model.User
	model.DB.Preload("Memo").Take(&user, ID)
	fmt.Println(user)
	c.JSON(http.StatusOK, response.Response{
		Code: http.StatusOK,
		Msg:  "success",
		Data: user.Memo,
	})
}

func MemoUpdate(c *gin.Context) {

}

func MemoDelete(c *gin.Context) {

}
