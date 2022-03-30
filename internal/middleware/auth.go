package middleware

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/TencentLexiang/auth-connector/internal/config"
	"github.com/TencentLexiang/auth-connector/pkg/webjwt"
	"github.com/gin-gonic/gin"
)

var AuthInfo struct {
	UserID string `json:"user_id"`
}

// Auth 判断是否携带登录态
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		loginURL := fmt.Sprintf("/page/login?referer=%s", url.QueryEscape(c.Request.RequestURI))
		if err != nil {
			c.Redirect(http.StatusFound, loginURL)
			return
		}
		err = getUserInfoFromToken(token)
		if err != nil {
			c.Redirect(http.StatusFound, loginURL)
			return
		}
		c.Next()
	}
}

func getUserInfoFromToken(token string) error {
	jwt := webjwt.NewJWT(config.Config.AuthJwtSecret)
	claims, err := jwt.ParseToken(token)
	if err != nil {
		return err
	}
	AuthInfo.UserID = claims.UserID
	return nil
}
