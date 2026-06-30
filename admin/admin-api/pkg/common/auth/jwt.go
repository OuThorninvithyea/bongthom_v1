package auth

import (

	// Community packages
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Claims holds the JWT payload — fields must match tbl_users + auth tables
type Claims struct {
	UserID       int64  `json:"user_id"`
	UserName     string `json:"user_name"`
	LoginSession string `json:"login_session"`
	RoleID       int    `json:"role_id"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed HS256 JWT valid for the given duration
func GenerateToken(userID int64, userName string, loginSession string, roleID int, secret string, d time.Duration) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(d)

	claims := Claims{
		UserID:       userID,
		UserName:     userName,
		LoginSession: loginSession,
		RoleID:       roleID,
		RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return signed, expiresAt, nil
}

// ValidateToken parses and validates a JWT, returning the claims on success
func ValidateToken(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
