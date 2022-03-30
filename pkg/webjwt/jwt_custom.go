package webjwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

// CreateTokenCustom 创建含有自定义字段的token
func (j *JWT) CreateTokenCustom(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseTokenCustom 解析含有自定义字段jwt token
func (j *JWT) ParseTokenCustom(tokenString string, claims jwt.Claims) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, fmt.Errorf("token validation error: %v", err)
			}
		}
	}

	if token.Valid {
		return claims, nil
	}

	return nil, TokenInvalid
}
