package router

import (
	"SimpleMemo/middleware"
	"SimpleMemo/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func UserRegister(c *gin.Context) {
	var regForm *model.User
	if err := c.ShouldBind(&regForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid register form",
			"err": err.Error(),
		})
		return
	}

	// Check if there's existed user
	var existed int64
	model.DB.Where("name = ?", regForm.Name).First(&model.User{}).Count(&existed)
	if existed != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "user already existed."})
		return
	}

	// Hash user password, algorithm bcrypt, cost 8
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(regForm.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "user data encryption failed."})
		return
	}
	model.DB.Create(&model.User{
		Name:      regForm.Name,
		Password:  string(encryptedPass),
		Email:     regForm.Email,
		CreatedAt: time.Now(),
		Memo:      nil,
	})
	c.JSON(http.StatusAccepted, gin.H{"msg": "user registration successful"})
}

func UserLogin(c *gin.Context) {
	// Bind login form
	var loginForm *model.UserLogin
	if err := c.ShouldBind(&loginForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid login form"})
		return
	}
	// Check if there's valid user
	// In order to prevent username scanning, all exceptions returns "invalid user or password" message.
	var user model.User
	if err := model.DB.Take(&user, "name = ?", loginForm.Name).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid user or password"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid user or password"})
		return
	}
	// Request JSON web token
	token, err := middleware.IssueToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "token issue failed",
			"err": err.Error(),
		})
		return
	}
	//token, _ := middleware.IssueToken(user.ID)
	c.JSON(http.StatusOK, gin.H{
		"msg":   "login successful",
		"token": token,
	})
}

func UserProfile(c *gin.Context) {
	ID, ok := c.Get("ID")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "no user info found in token"})
	}
	var user *model.User
	if err := model.DB.Take(&user, "ID = ?", ID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": gin.H{
			"ID":        user.ID,
			"Name":      user.Name,
			"Email":     user.Email,
			"CreatedAt": user.CreatedAt,
		},
	})
}

func UserUpdate(c *gin.Context) {
	// Check if user is valid
	ID, ok := c.Get("ID")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "no user info found in token"})
	}
	var user *model.User
	if err := model.DB.Take(&user, "ID = ?", ID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	// Check data for update
	if name, ok := c.GetPostForm("Name"); ok {
		var existed int64
		model.DB.Where("name = ?", name).First(&model.User{}).Count(&existed)
		if existed != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "username already existed, please use another one"})
			return
		}
		user.Name = name
	}
	if email, ok := c.GetPostForm("Email"); ok {
		var existed int64
		model.DB.Where("email = ?", email).First(&model.User{}).Count(&existed)
		if existed != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "email already existed, please use another one"})
			return
		}
		user.Email = email
	}
	if pass, ok := c.GetPostForm("Password"); ok {
		if prevPass, ok := c.GetPostForm("PrevPassword"); ok {
			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(prevPass)); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "changing password without correct password"})
				return
			} else {
				encryptedPass, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"msg": "user data encryption failed."})
					return
				}
				user.Password = string(encryptedPass)
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "changing password without previous password"})
			return
		}
	}
	// Save record to database, return status ok
	model.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{
		"msg": "data update successful",
		"info": gin.H{
			"ID":        user.ID,
			"Name":      user.Name,
			"Email":     user.Email,
			"CreatedAt": user.CreatedAt,
		},
	})
}
