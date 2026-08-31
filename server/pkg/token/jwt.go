package token

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID  uint   `json:"user_id"`
	Role    string `json:"role"`
	JTI     string `json:"jti"`
	jwt.RegisteredClaims
}

type AccessToken struct {
	Token     string
	JTI       string
	ExpiresAt time.Time
}

func GenerateAccessToken(userID uint, isAdmin bool, secret string, ttl time.Duration) (*AccessToken, error) {
	jti := uuid.NewString()
	now := time.Now()
	exp := now.Add(ttl)
	role := "user"
	if isAdmin {
		role = "admin"
	}
	claims := Claims{
		UserID: userID,
		Role:   role,
		JTI:    jti,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			ID:        jti,
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	return &AccessToken{Token: signed, JTI: jti, ExpiresAt: exp}, nil
}

func ParseAccessToken(tokenStr, secret string) (*Claims, error) {
	claims := &Claims{}
	t, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func RandomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
