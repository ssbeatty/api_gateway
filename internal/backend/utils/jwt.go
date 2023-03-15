package utils

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTTokenGen struct {
	privateKey *rsa.PrivateKey
	issuer     string
	exp        time.Duration
}

// NewJWTTokenGen creates a JWTTokenGen.
func NewJWTTokenGen(issuer string, privateKey *rsa.PrivateKey) *JWTTokenGen {
	return &JWTTokenGen{
		issuer:     issuer,
		privateKey: privateKey,
	}
}

// GenerateToken generates token
func (t *JWTTokenGen) GenerateToken(userName string, expire time.Duration) (string, error) {
	nowSec := time.Now().Unix()

	claims := jwt.StandardClaims{
		Issuer:    t.issuer,
		IssuedAt:  nowSec,
		ExpiresAt: nowSec + int64(expire.Seconds()),
		Subject:   userName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	return token.SignedString(t.privateKey)
}
