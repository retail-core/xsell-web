package session

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Sub        string `json:"sub"`
	Role       string `json:"role"`
	IsVerified bool   `json:"is_verified"`
	TenantID   string `json:"tenant_id,omitempty"`
}

func DecodeAndVerify(tokenStr, secret string) (*Claims, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	result := &Claims{}
	if sub, ok := claims["sub"].(string); ok {
		result.Sub = sub
	}
	if role, ok := claims["role"].(string); ok {
		result.Role = role
	}
	if verified, ok := claims["is_verified"].(bool); ok {
		result.IsVerified = verified
	}
	if tenant, ok := claims["tenant_id"].(string); ok {
		result.TenantID = tenant
	}

	return result, nil
}