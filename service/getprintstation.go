package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"strings"
)

type stationListRes struct {
	Errmsg string `json:"ErrorMessage"`
	Result []*station `json:"Result"`
}

type station struct {
	DevSN int32 `json:"dwDevSN"`
	CtrlType int32 `json:"dwCtrlType"`
	PropertyCode int32 `json:"dwProperty"`
	StatusCode int32 `json:"dwStatus"`
	OpenTime int32 `json:"dwOpenTime"`
	CloseTime int32 `json:"dwCloseTime"`
	Property string `json:"szProperty"`
	Status string `json:"szStatus"`
	Form string `json:"szForm"`
	Name string `json:"szName"`
	Position string `json:"szPosition"`
	Tel string `json:"szTel"`
	RoomName string `json:"szRoomName"`
	StationInfo string `json:"szStatInfo"`
	Memo string `json:"szMemo"`
}

func requestStation(sessionID string) ([]*station, error) {
	url := "http://print.intl.zju.edu.cn/Service.asmx/GetDevices"
	payload := strings.NewReader(fmt.Sprintf("{\"bstrSessionID\": \"%s\"}", sessionID))
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,la;q=0.7")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Length", "39")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", "SessionID=636788737154315827")
	req.Header.Add("Host", "print.intl.zju.edu.cn")
	req.Header.Add("Origin", "http://print.intl.zju.edu.cn")
	req.Header.Add("Referer", "http://print.intl.zju.edu.cn/")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	data := &stationListRes{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	if data.Errmsg != "" {
		return nil, errors.New(data.Errmsg)
	}

	for _, item := range data.Result {
		docStatus, err := goquery.NewDocumentFromReader(strings.NewReader(item.Status))
		if err != nil {
			continue
		}
		statusText := docStatus.Find("span").Text()


		docStationInfo, err := goquery.NewDocumentFromReader(strings.NewReader(item.StationInfo))
		if err != nil {
			continue
		}
		stationInfoText := docStationInfo.Find("span").Text()

		item.Status = statusText
		item.StationInfo = stationInfoText
	}

	return data.Result, nil
}

func GetPrintStation() ([]*station, error) {
	session, err := getSession()
	if err != nil {
		return nil, err
	}
	//if err := simiLogin(session, username, password); err != nil {
	//	return nil, err
	//}
	return requestStation(session)
}