# auth-connector

## 概述
此项目为接入[腾讯乐享](https://lexiangla.com)企业登录示例代码，提供简单的登录授权流程示例。

此项目代码请勿直接用于企业生产环境，仅作参考

## 本地使用
安装docker

docker build -t auth-connector .

docker run -d -p 80:9000 auth-connector

浏览器访问：http://127.0.0.1/page/lxauth?redirect_uri=https://lexiangla.com/suites/auth-callback&state=123 
跳转到 https://lexiangla.com/suites/auth-callback?code=53fe3239-cb78-4c40-8966-e8b1d52eb017&state=123
则表示服务正常

## 介绍

此项目提供3个cgi:
- /page/lxauth 
- /page/login
- /api/userinfo

### /page/lxauth
授权页，用于生成一次性临时授权码并回跳至乐享授权回调页，临时授权码应当包含当前用户态信息，
因此此cgi依赖用户登录态cookie，若cookie不存在或登录态过期，则跳转至`/page/login`页面进行登录

### /page/login
登陆页，对应企业OA系统的登录页面，用户输入登录凭证后生成登录态cookie。此cgi无需配置在乐享表单中。

### /api/userinfo
获取用户信息接口，由乐享服务端使用上述临时授权码调用此接口获取用户信息

