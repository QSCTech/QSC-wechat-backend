package service

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

func simiLogin(session, username, password string) error {
	url := "http://print.intl.zju.edu.cn/Service.asmx"
	payload := strings.NewReader(fmt.Sprintf("<?xml version=\"1.0\" encoding=\"utf-8\"?><soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:soapenc=\"http://schemas.xmlsoap.org/soap/encoding/\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" ><soap:Body><Login xmlns=\"http://tempuri.org/\"><bstrSessionID>%s</bstrSessionID><bstrUserName>%s</bstrUserName><bstrPassword>%s</bstrPassword></Login></soap:Body></soap:Envelope>", session, username, password))
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Postman-Token", "d8841451-b0c1-4879-9143-7e5b583e0179")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}
	resText := doc.Find("LoginResult").Text()
	result := strings.Split(resText, ",")[0]
	if result != "ok" {
		return errors.New("Password not correct")
	}
	return nil
}

func getSession() (string, error) {
	url := "http://print.intl.zju.edu.cn/Service.asmx"
	payload := strings.NewReader("<?xml version=\"1.0\" encoding=\"utf-8\"?><soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:soapenc=\"http://schemas.xmlsoap.org/soap/encoding/\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" ><soap:Body><InitSession xmlns=\"http://tempuri.org/\"><bstrPCName></bstrPCName></InitSession></soap:Body></soap:Envelope>")
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("cache-control", "no-cache")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}
	sessionText := doc.Find("InitSessionResult").Text()
	session := strings.Split(sessionText, ",")[1]
	return session, nil
}

func CheckPrintPassword(username, password string) error {
	session, err := getSession()
	if err != nil {
		return err
	}
	return simiLogin(session, username, password)
}