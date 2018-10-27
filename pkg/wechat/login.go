package wechat

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

type session struct {
	OpenID string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID string `json:"unionid"`
}

func Code2Session(js_code string) *session {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		viper.GetString("wechat.AppID"), viper.GetString("wechat.AppSecret"), js_code)
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("cache-control", "no-cache")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(res)
	fmt.Println(string(body))
}