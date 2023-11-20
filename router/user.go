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
	var regForm *model.RegForm
	if err := c.ShouldBind(&regForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid register form",
			"err": err,
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
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(regForm.Password), 8)
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
	var loginForm *model.LoginForm
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
			"err": err,
		})
		return
	}
	//token, _ := middleware.IssueToken(user.ID)
	c.JSON(http.StatusAccepted, gin.H{
		"msg":   "login successful",
		"token": token,
	})
}
