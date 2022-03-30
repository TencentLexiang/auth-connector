package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/TencentLexiang/auth-connector/internal/config"
	"github.com/TencentLexiang/auth-connector/internal/middleware"
	"github.com/TencentLexiang/auth-connector/pkg/webjwt"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid" //nolint:goimports
)

func init() {
	// suppose to load config from yaml
	config.Config.AuthJwtSecret = "example"
}

func main() {
	r := gin.Default()
	useAuth := r.Group("/")
	useAuth.Use(middleware.Auth())
	{
		useAuth.GET("/page/lxauth", lxauth)
	}
	r.GET("/page/login", login)
	r.POST("/api/userinfo", userinfo)
	err := r.Run("0.0.0.0:9000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		panic("server fail")
	}
}

func lxauth(c *gin.Context) {
	state := c.Query("state")
	if state == "" {
		c.JSON(400, gin.H{"message": "`state` is required"})
		return
	}
	redirectURI := c.Query("redirect_uri")
	r, err := url.QueryUnescape(redirectURI)
	if redirectURI == "" || err != nil {
		c.JSON(400, gin.H{"message": "`redirect_uri` error"})
		return
	}

	code := uuid.NewV4()
	fmt.Printf("state:%s, code:%s", state, code)
	// todo: use redis to save kv, key for code and value for userid，for userinfo api
	// redis command: SETEX uuid 180 {middleware.AuthInfo.UserID}
	c.Redirect(http.StatusFound, fmt.Sprintf("%s?code=%s&state=%s", r, code, state))
}

func login(c *gin.Context) {
	referer := c.Query("referer")
	jwt := webjwt.NewJWT(config.Config.AuthJwtSecret)
	token, err := jwt.CreateToken(webjwt.CustomClaims{
		UserID: "LX001",
	})
	if err != nil {
		c.JSON(400, gin.H{"message": "create jwt error"})
		return
	}
	c.SetCookie("token", token, 3600, "/", "/", false, false)
	r, err := url.QueryUnescape(referer)
	if err != nil {
		c.JSON(400, gin.H{"message": "referer error"})
		return
	}
	c.Redirect(http.StatusFound, r)
}

func userinfo(c *gin.Context) {
	code := c.PostForm("code")
	fmt.Print("code", code)
	// todo: userid should be gotten from redis
	// redis command: GET uuid && del uuid
	userid := "LX001"
	c.JSON(200, gin.H{"id": userid})
}
