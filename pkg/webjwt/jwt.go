// Package webjwt web jwt library.
package webjwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWT jwt签名
type JWT struct {
	SigningKey []byte
}

// JwtConfig jwt config.
type JwtConfig struct {
	Key string `yaml:"jwt_secret"`
}

// TokenLeftTime token expire hours.
const TokenLeftTime = 60 * 24 // 过期小时

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that is not even a token")
	TokenInvalid     = errors.New("token invalid")
)

// CustomClaims 载荷
type CustomClaims struct {
	jwt.StandardClaims
	UserID string `json:"user_id,omitempty"`
}

// NewJWT 创建一个jwt实例
func NewJWT(key string) *JWT {
	return &JWT{
		SigningKey: []byte(key),
	}
}

// GetSignKey 获取SignKey
func (j *JWT) GetSignKey() string {
	return string(j.SigningKey)
}

// SetSignKey 设置signkey
func (j *JWT) SetSignKey(key string) {
	j.SigningKey = []byte(key)
}

// CreateToken 创建token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析jwt token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
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

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, TokenInvalid
}

// RefreshToken 刷新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(TokenLeftTime * time.Hour).Unix()
		return j.CreateToken(*claims)
	}

	return "", TokenInvalid
}
