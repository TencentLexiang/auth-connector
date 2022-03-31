package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/TencentLexiang/auth-connector/internal/config"
	"github.com/TencentLexiang/auth-connector/internal/middleware"
	"github.com/TencentLexiang/auth-connector/pkg/webjwt"
	"github.com/TencentLexiang/auth-connector/pkg/workwechat"
	_ "github.com/TencentLexiang/auth-connector/pkg/workwechat"
	"github.com/gin-gonic/gin"
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
	r.GET("/page/auth-callback", cb)
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

	// todo: code 应该使用`uuid.NewV4()`生成，为了演示企业微信授权流程，直接使用企业微信用户ID，生产环境严禁使用，否则存在极大的安全隐患
	code := middleware.AuthInfo.UserID
	fmt.Printf("state:%s, code:%s", state, code)
	// todo: use redis to save kv, key for code and value for userid，for userinfo api
	// redis command: SETEX uuid 180 {middleware.AuthInfo.UserID}
	c.Redirect(http.StatusFound, fmt.Sprintf("%s?code=%s&state=%s", r, code, state))
}

func login(c *gin.Context) {
	referer := c.Query("referer")
	appid := workwechat.Config.CorpID
	cbURL := c.Request.Host + "/page/auth-callback"
	u := fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&"+
		"redirect_uri=%s&response_type=code&scope=snsapi_base&state=%s#wechat_redirect",
		appid, url.QueryEscape(cbURL), url.QueryEscape(referer))
	c.Redirect(http.StatusFound, u)
}

func cb(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	if code == "" {
		c.JSON(400, gin.H{"message": "`code` error"})
		return
	}
	// 从企业微信获取用户身份
	userid := workwechat.GetUserInfo(code)

	// 把用户身份注入站点cookie，避免频繁调用企业微信的授权
	jwt := webjwt.NewJWT(config.Config.AuthJwtSecret)
	token, err := jwt.CreateToken(webjwt.CustomClaims{
		UserID: userid,
	})
	if err != nil {
		c.JSON(400, gin.H{"message": "create jwt error"})
		return
	}
	c.SetCookie("token", token, 3600, "/", "/", false, false)
	r, err := url.QueryUnescape(state)
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
	userid := code
	c.JSON(200, gin.H{"id": userid})
}
