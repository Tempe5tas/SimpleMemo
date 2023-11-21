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
			Msg:  "invalid time format, should be YYYY-MM-DD HH:MM:SS",
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
	// Get memo ID
	ID, ok := c.GetPostForm("ID")
	if !ok {
		c.JSON(http.StatusBadRequest, response.Response{
			Code: http.StatusBadRequest,
			Msg:  "no memo ID provided",
		})
		return
	}
	// Verify ownership
	var memo model.Memo
	model.DB.Take(&memo, "ID = ?", ID)
	UID, ok2 := c.Get("ID")
	if !ok2 || memo.UID != UID {
		c.JSON(http.StatusUnauthorized, response.Response{
			Code: http.StatusUnauthorized,
			Msg:  "permission denied: not memo owner",
		})
		return
	}
	// Check request body for update
	if Title, ok := c.GetPostForm("Title"); ok {
		memo.Title = Title
	}
	if Content, ok := c.GetPostForm("Content"); ok {
		memo.Content = Content
	}
	if Time, ok := c.GetPostForm("Time"); ok {
		memoTime, err := time.Parse("2006-01-02 15:04:05", Time)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Response{
				Code: http.StatusBadRequest,
				Msg:  "invalid time format, should be YYYY-MM-DD HH:MM:SS",
				Data: err.Error(),
			})
			return
		}
		memo.Time = memoTime
	}
	if Status, ok := c.GetPostForm("Status"); ok {
		if Status == "true" {
			memo.Status = true
		} else if Status == "false" {
			memo.Status = false
		} else {
			c.JSON(http.StatusBadRequest, response.Response{
				Code: http.StatusBadRequest,
				Msg:  "invalid memo format: not true/false in \"Status\"",
			})
			return
		}
	}
	// All variables retrieved, start database record saving
	model.DB.Save(&memo)
	c.JSON(http.StatusOK, response.Response{
		Code: http.StatusOK,
		Msg:  "memo update successful",
		Data: memo,
	})
}

func MemoDelete(c *gin.Context) {
	// Get memo ID
	ID, ok := c.GetPostForm("ID")
	if !ok {
		c.JSON(http.StatusBadRequest, response.Response{
			Code: http.StatusBadRequest,
			Msg:  "no memo ID provided",
		})
		return
	}
	// Verify ownership
	var memo model.Memo
	model.DB.Take(&memo, "ID = ?", ID)
	UID, ok2 := c.Get("ID")
	if !ok2 || memo.UID != UID {
		c.JSON(http.StatusUnauthorized, response.Response{
			Code: http.StatusUnauthorized,
			Msg:  "permission denied: not memo owner",
		})
		return
	}
	// Delete memo record from database
	model.DB.Delete(&memo)
	c.JSON(http.StatusOK, response.Response{
		Code: http.StatusOK,
		Msg:  "memo delete successful",
	})
}
