package utils

import (
	"api_gateway/internal/backend/config"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type JwtComponent struct{}

func GetClaimsByToken(authHeader string) (jwt.MapClaims, error) {
	if authHeader == "" {
		return nil, ErrAuthToken
	}
	jwtConfig := config.DefaultConfig.Jwt
	if token, err := ValidateToken(authHeader[len(jwtConfig.BearerSchema):]); err != nil {
		return nil, ErrAuthToken
	} else {
		if claims, ok := token.Claims.(jwt.MapClaims); !ok {
			return nil, ErrAuthToken
		} else {
			if token.Valid {
				return claims, nil
			} else {
				return nil, ErrAuthToken
			}
		}
	}
}

// GenerateToken generates token
func GenerateToken(userId int, userName string) (string, int64, error) {
	jwtConfig := config.DefaultConfig.Jwt
	exp := time.Now().Add(time.Minute * jwtConfig.JwtExpMin).Unix()
	claims := jwt.MapClaims{
		"exp":       exp,
		"iat":       time.Now().Unix(),
		"user_id":   userId,
		"user_name": userName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(jwtConfig.JwtSecret))
	return t, exp, err
}

// ValidateToken validate the given token
func ValidateToken(token string) (*jwt.Token, error) {
	jwtConfig := config.DefaultConfig.Jwt
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			//nil secret key
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtConfig.JwtSecret), nil
	})
}

func JWT() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 获取token
		token := context.GetHeader("Token")

		if token == "" {
			context.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "没有携带token",
				"data": "",
			})
			context.Abort()
			return
		} else {
			Claims, err := GetClaimsByToken(token)
			if err != nil {
				context.JSON(http.StatusOK, gin.H{
					"code": http.StatusUnauthorized,
					"msg":  "token验证失败",
					"data": "",
				})
				context.Abort()
				return
			} else if time.Now().Unix() > Claims["exp"].(int64) {
				context.JSON(http.StatusOK, gin.H{
					"code": http.StatusUnauthorized,
					"msg":  "token已过期",
					"data": "",
				})
				context.Abort()
				return
			}
		}
	}
}
