package utils

import (
	"api_gateway/internal/backend/config"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cast"
	"time"
)

type JwtComponent struct{}

func (j *JwtComponent) GetUserIdByToken(ctx context.Context, authHeader string) (userId int64, err error) {
	if authHeader == "" {
		return 0, ErrAuthToken
	}
	jwtConfig := config.DefaultConfig.Jwt
	if token, err := ValidateToken(authHeader[len(jwtConfig.BearerSchema):]); err != nil {
		return 0, ErrAuthToken
	} else {
		if claims, ok := token.Claims.(jwt.MapClaims); !ok {
			return 0, ErrAuthToken
		} else {
			if token.Valid {
				userId := cast.ToInt64(claims["user_id"])
				if userId == 0 {
					return 0, ErrAuthToken
				}

				return userId, nil
			} else {
				return 0, ErrAuthToken
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
