package middleware

import (
	response "SimpleMemo/serializer"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

type MyClaims struct {
	ID uint
	jwt.StandardClaims
}

// JWT token issuing and parsing function.

func IssueToken(ID uint) (string, error) {
	claim := MyClaims{
		ID,
		jwt.StandardClaims{
			Issuer:    viper.GetString("JWT.Issuer"),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Duration(viper.GetInt("JWT.Expire")) * time.Second).Unix(),
		}}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(viper.GetString("JWT.Secret")))
}

func ParseToken(str string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(str, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("JWT.Secret")), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// Real middleware for user token validating.

func ValidateToken(c *gin.Context) {
	tokenStr := c.Request.Header.Get("token")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, response.Response{
			Code: http.StatusUnauthorized,
			Msg:  "token missing, please login",
		})
		c.Abort()
		return
	}
	claim, err := ParseToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Response{
			Code: http.StatusUnauthorized,
			Msg:  err.Error(),
		})
		c.Abort()
		return
	}
	if viper.GetString("JWT.Issuer") != claim.Issuer {
		c.JSON(http.StatusUnauthorized, response.Response{
			Code: http.StatusUnauthorized,
			Msg:  "invalid token issuer, validation aborted",
		})
		c.Abort()
		return
	}
	//if time.Now().Unix() < claim.NotBefore {
	//	c.JSON(http.StatusUnauthorized, gin.H{"msg": "Token not activated yet, please retry."})
	//	c.Abort()
	//	return
	//}
	//if time.Now().Unix() > claim.ExpiresAt {
	//	c.JSON(http.StatusUnauthorized, gin.H{"msg": "Token expired, please login."})
	//	c.Abort()
	//	return
	//}
	c.Set("ID", claim.ID)
	c.Next()
}
