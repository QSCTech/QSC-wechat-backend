package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type jobListRes struct {
	Errmsg string `json:"ErrorMessage"`
	Result []*printJob `json:"Result"`
}

type printJob struct {
	JobId int32 `json:"dwJobId"`
	Pages int32 `json:"dwPages"`
	Copies int32 `json:"dwCopies"`
	Color int32 `json:"dwColor"`
	OutMode int64 `json:"dwOutMode"`
	Time string `json:"szDateTime"`
	Form string `json:"szForm"`
	Property string `json:"szProperty"`
	Document string `json:"szDocument"`
	JobFileName string `json:"szJobFileName"`
	Memo string `json:"szMemo"`
}

func requestGetJob(sessionID string) ([]*printJob, error) {
	url := "http://print.intl.zju.edu.cn/Service.asmx/GetPrintJob"
	payload := strings.NewReader(fmt.Sprintf("{\"bstrSessionID\": \"%s\"}", sessionID))
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,la;q=0.7")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Length", "39")
	req.Header.Add("Content-Type", "application/json")
	//req.Header.Add("Cookie", "Hm_lvt_309754640d9cba6fb911afd186cdd934=1541764687,1541775184,1542675357,1542854299; SESSc14f986e6881fe9ccdf82d6eb620f77e=L06EdhAYxbmqw_k0Z9APXr_U_-G3y0nUP0GRCiV8gkQ; ASP.NET_SessionId=5a5w0f555irvv0450brqhh55; SessionID=636788698645024189; Logined=1; LoginInfo=360345%#3170111705%#%#%#%#%#Z09030742%#ç½æå¿%#8535%#0%$")
	req.Header.Add("Host", "print.intl.zju.edu.cn")
	req.Header.Add("Origin", "http://print.intl.zju.edu.cn")
	req.Header.Add("Referer", "http://print.intl.zju.edu.cn/")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	data := &jobListRes{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	if data.Errmsg != "" {
		return nil, errors.New(data.Errmsg)
	}

	return data.Result, nil
}

func GetPrintJob(username, password string) ([]*printJob, error) {
	session, err := getSession()
	if err != nil {
		return nil, err
	}
	if err := simiLogin(session, username, password); err != nil {
		return nil, err
	}

	return requestGetJob(session)
}