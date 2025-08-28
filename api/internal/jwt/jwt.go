package jwtpkg

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/happilymarrieddad/order-management-v3/api/types"
)

// For production, this should be loaded from a secure configuration management system or environment variables.
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// CustomClaims includes custom data for the JWT, embedding standard claims.
type CustomClaims struct {
	UserID    int64       `json:"userId"`
	Email     string      `json:"email"`
	CompanyID int64       `json:"companyId"`
	Roles     types.Roles `json:"roles"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT for a given user.
func GenerateToken(user *types.User) (string, error) {
	// Set custom claims
	claims := &CustomClaims{
		UserID:    user.ID,
		Email:     user.Email,
		CompanyID: user.CompanyID,
		Roles:     user.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)), // Token expires in 3 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "order-management-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and return it as a string.
	return token.SignedString(jwtSecret)
}

// ValidateToken parses and validates a token string, returning the custom claims if valid.
func ValidateToken(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}
