package handlers

import (
	"ginvueblog/config"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"fmt"
)


// qq回调逻辑
func  Index(c *gin.Context) {
	html:=`<!DOCTYPE html>
	<html lang="en">
	<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="ie=edge">
	<title>飞书Go技术分享</title>
	</head>
	<body>
		飞书Go技术分享
		<br/>
		飞书Go技术分享
	</body>
`
	c.Status(200)
	tmpl, err := template.New("htmlTest").Parse(html)
	if err != nil {
		fmt.Printf("parsing: %s", err)
	}
	tmpl.Execute(c.Writer, map[string]interface{}{
		"test1": "test",
	})

}



// qq回调逻辑
func CallBack(c *gin.Context) {

	code, _ := c.GetPostForm("code")

	params := url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("client_id", config.AppId)
	params.Add("client_secret", config.AppKey)
	params.Add("code", code)
	str := fmt.Sprintf("%s&redirect_uri=%s", params.Encode(), config.RedirectURI)
	loginURL := fmt.Sprintf("%s?%s", "https://graph.qq.com/oauth2.0/token", str)

	response, err := http.Get(loginURL)
	defer response.Body.Close()

	if err != nil {
		//w.Write([]byte(err.Error()))
		c.JSON(http.StatusBadRequest, err.Error())
	}

	bs, _ := ioutil.ReadAll(response.Body)
	body := string(bs)

	resultMap := convertToMap(body)

	info := &config.PrivateInfo{}
	info.AccessToken = resultMap["access_token"]
	info.RefreshToken = resultMap["refresh_token"]
	info.ExpiresIn = resultMap["expires_in"]

	GetOpenId(info, c)

	//
}

// 1. Get Authorization Code
func GetAuthCode(c *gin.Context) {
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", config.AppId)
	params.Add("state", "test")
	str := fmt.Sprintf("%s&redirect_uri=%s", params.Encode(), config.RedirectURI)
	loginURL := fmt.Sprintf("%s?%s", "https://graph.qq.com/oauth2.0/authorize", str)

	c.Redirect(http.StatusFound, loginURL)
	//http.Redirect(w, r, loginURL, http.StatusFound)
}

// 2. Get Access Token
func GetToken(c *gin.Context) {
	// use CallBack
}

// 3. Get OpenId
func GetOpenId(info *config.PrivateInfo, c *gin.Context) {
	resp, err := http.Get(fmt.Sprintf("%s?access_token=%s", "https://graph.qq.com/oauth2.0/me", info.AccessToken))
	if err != nil {
		//w.Write([]byte(err.Error()))
		c.JSON(http.StatusBadRequest, err.Error())
	}
	defer resp.Body.Close()

	bs, _ := ioutil.ReadAll(resp.Body)
	body := string(bs)
	info.OpenId = body[45:77]

	GetUserInfo(info, c)
}

// 4. Get User info
func  GetUserInfo(info *config.PrivateInfo, c *gin.Context) {
	params := url.Values{}
	params.Add("access_token", info.AccessToken)
	params.Add("openid", info.OpenId)
	params.Add("oauth_consumer_key", config.AppId)

	uri := fmt.Sprintf("https://graph.qq.com/user/get_user_info?%s", params.Encode())
	resp, err := http.Get(uri)
	if err != nil {
		//w.Write([]byte(err.Error()))
		c.JSON(http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusOK, resp.Body)
}

func convertToMap(str string) map[string]string {
	var resultMap = make(map[string]string)
	values := strings.Split(str, "&")
	for _, value := range values {
		vs := strings.Split(value, "=")
		resultMap[vs[0]] = vs[1]
	}
	return resultMap
}