// Package token: สร้าง/ตรวจ JWT + Claims — auth-specific ของ user-service.
package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID        string `json:"uid"`
	Username      string `json:"username"`
	FullName      string `json:"full_name"`
	Email         string `json:"email"`
	Role          string `json:"role"`
	Provider      string `json:"provider"`       // login มาทางไหน (นิ่ง ฝังได้)
	EmailVerified bool   `json:"email_verified"` // ยืนยันอีเมลแล้วไหม
	jwt.RegisteredClaims
}

func Generate(c Claims, secret string, ttl time.Duration) (string, error) {
	now := time.Now()
	c.RegisteredClaims = jwt.RegisteredClaims{
		Subject:   c.UserID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(secret))
}

func Parse(tokenStr, secret string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}
