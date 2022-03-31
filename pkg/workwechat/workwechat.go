package workwechat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var Config struct {
	CorpID     string
	CorpSecret string
}

func init() {
	// todo: init config
	Config.CorpID = ""
	Config.CorpSecret = ""
	if Config.CorpID == "" || Config.CorpSecret == "" {
		panic("请先完成企业微信配置")
	}
}

func GetAccessToken() string {
	u := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s",
		Config.CorpID, Config.CorpSecret)
	rs, err := http.Get(u)
	if err != nil {
		log.Print("请求企业微信gettoken接口失败")
	}
	defer rs.Body.Close()
	rsp, _ := ioutil.ReadAll(rs.Body)

	var tokenStruct struct {
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
		ExpiresIn   string `json:"expires_in"`
	}
	err = json.Unmarshal(rsp, &tokenStruct)
	if err != nil {
		log.Print("json Unmarshal error")
	}
	return tokenStruct.AccessToken
}

func GetUserInfo(code string) string {
	u := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=%s&code=%s",
		GetAccessToken(), code)
	rs, err := http.Get(u)
	if err != nil {
		log.Print("请求企业微信getuserinfo接口失败")
	}
	defer rs.Body.Close()
	rsp, _ := ioutil.ReadAll(rs.Body)

	var UserInfo struct {
		ErrCode  int    `json:"errcode"`
		ErrMsg   string `json:"errmsg"`
		UserID   string `json:"UserId"`
		DeviceID string `json:"DeviceId"`
	}
	err = json.Unmarshal(rsp, &UserInfo)
	if err != nil {
		log.Print("json Unmarshal error")
	}
	return UserInfo.UserID
}
