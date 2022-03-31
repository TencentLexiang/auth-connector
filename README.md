# auth-connector

## 概述
此项目为接入[腾讯乐享](https://lexiangla.com)企业登录示例代码，提供简单的登录授权流程示例。

**此分支为对接企业微信自建应用授权示例，需要提前注册企业微信并创建自建应用** 

此项目代码请勿直接用于企业生产环境，仅作参考

## 本地使用
在 pkg/workwechat/workwechat.go 下配置企业微信corpID以及自建账号corpSecret

安装docker

docker build -t auth-connector .

docker run -d -p 80:9000 auth-connector

浏览器访问：http://127.0.0.1/page/lxauth?redirect_uri=https://lexiangla.com/suites/auth-callback&state=123 
跳转到 https://open.weixin.qq.com/ 则表示服务正常

## 介绍

此项目提供3个cgi:
- /page/lxauth 
- /page/login
- /page/auth-callback
- /api/userinfo

### /page/lxauth
授权页，用于生成一次性临时授权码并回跳至乐享授权回调页，临时授权码应当包含当前用户态信息，
因此此cgi依赖用户登录态cookie，若cookie不存在或登录态过期，则跳转至`/page/login`页面进行登录

### /page/login
登陆页，在对接企业微信自建账号流程中，用户登录态依赖企业微信授权，因此该页面会直接重定向至企业微信授权页。

### /page/auth-callback
企业微信授权回调页，使用企业微信的一次性临时授权码请求企业微信接口，获取用户信息，生成登录态cookie，
并重定向至/page/lxauth继续完成乐享的授权

### /api/userinfo
获取用户信息接口，由乐享服务端使用上述临时授权码调用此接口获取用户信息

