package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type delJobRes struct {
	Errmsg string `json:"ErrorMessage"`
}
func requestDelJob(sessionID string, jobId int32) (error) {
	url := "http://print.intl.zju.edu.cn/Service.asmx/SetPrintJob"

	payload := strings.NewReader(fmt.Sprintf("{\"bstrSessionID\": \"%s\",\"bstrOP\":\"DEL\",\"bstrJobID\":\"%d\"}", sessionID, jobId))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,la;q=0.7")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Length", "75")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Host", "print.intl.zju.edu.cn")
	req.Header.Add("Origin", "http://print.intl.zju.edu.cn")
	req.Header.Add("Referer", "http://print.intl.zju.edu.cn/")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	data := &delJobRes{}
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	if data.Errmsg != "" {
		return errors.New(data.Errmsg)
	}
	return nil
}


func DeletePrintJob(username, password string, jobId int32) error {
	session, err := getSession()
	if err != nil {
		return err
	}
	if err := simiLogin(session, username, password); err != nil {
		return err
	}
	requestGetJob(session)
	return requestDelJob(session, jobId)
}