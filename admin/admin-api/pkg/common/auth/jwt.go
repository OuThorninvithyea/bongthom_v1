package auth

import (
	// Community packages
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims holds the JWT payload and the database-backed login session.
type Claims struct {
	UserID               int64  `json:"user_id"`
	UserName             string `json:"user_name"`
	RoleID               int    `json:"role_id"`
	LoginSession         string `json:"login_session"`
	jwt.RegisteredClaims        // exp, iat
}

// GenerateToken creates a signed HS256 JWT valid for the given duration
// and includes the session value stored on the user row.
func GenerateToken(userID int64, userName string, roleID int, loginSession string, secret string, d time.Duration) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(d)

	claims := Claims{
		UserID:       userID,
		UserName:     userName,
		RoleID:       roleID,
		LoginSession: loginSession,
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
