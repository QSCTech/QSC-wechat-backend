package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/minio/minio-go"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

func simiLogin(session, username, password string) error {
	url := "http://print.intl.zju.edu.cn/Service.asmx"
	payload := strings.NewReader(fmt.Sprintf("<?xml version=\"1.0\" encoding=\"utf-8\"?><soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:soapenc=\"http://schemas.xmlsoap.org/soap/encoding/\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" ><soap:Body><Login xmlns=\"http://tempuri.org/\"><bstrSessionID>%s</bstrSessionID><bstrUserName>%s</bstrUserName><bstrPassword>%s</bstrPassword></Login></soap:Body></soap:Envelope>", session, username, password))
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("cache-control", "no-cache")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}
	resText := doc.Find("LoginResult").Text()
	result := strings.Split(resText, ",")[0]
	if result != "ok" {
		return errors.New("password incorrect")
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

func newUploadRequest(url string, params map[string]string, file *multipart.FileHeader) error {
	f, err := file.Open() // 打开文件句柄
	if err != nil {
		return err
	}
	defer f.Close()

	body := &bytes.Buffer{} // 初始化body参数
	writer := multipart.NewWriter(body) // 实例化multipart
	part, err := writer.CreateFormFile("file", file.Filename) // 创建multipart 文件字段
	if err != nil {
		return err
	}

	_, err = io.Copy(part, f) // 写入文件数据到multipart
	if err != nil {
		return err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val) // 写入body中额外参数，比如七牛上传时需要提供token
	}
	err = writer.Close()

	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("POST", url, body) // 新建请求
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType()) // 设置请求头,!!!非常重要，否则远端无法识别请求

	res, err := client.Do(req)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func Print(username, password string, fields map[string]string, file *multipart.FileHeader) error {
	session, err := getSession()
	if err != nil {
		return err
	}
	if err := simiLogin(session, username, password); err != nil {
		return err
	}

	return newUploadRequest(fmt.Sprintf("http://print.intl.zju.edu.cn/upload.aspx?sid=%s", session),
		fields, file)
}



func minioUploadRequest(url string, params map[string]string, obj *minio.Object, fileName string) error {
	body := &bytes.Buffer{} // 初始化body参数
	writer := multipart.NewWriter(body) // 实例化multipart
	part, err := writer.CreateFormFile("file", fileName) // 创建multipart 文件字段
	if err != nil {
		return err
	}
	_, err = io.Copy(part, obj) // 写入文件数据到multipart
	if err != nil {
		return err
	}
	for key, val := range params {
		_ = writer.WriteField(key, val) // 写入body中额外参数，比如七牛上传时需要提供token
	}
	err = writer.Close()
	if err != nil {
		return err
	}
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("POST", url, body) // 新建请求
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType()) // 设置请求头,!!!非常重要，否则远端无法识别请求

	res, err := client.Do(req)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func PrintFromMinio(username, password string, fields map[string]string, object *minio.Object, fileName string) error {
	session, err := getSession()
	if err != nil {
		return err
	}
	if err := simiLogin(session, username, password); err != nil {
		return err
	}

	return minioUploadRequest(fmt.Sprintf("http://print.intl.zju.edu.cn/upload.aspx?sid=%s", session),
		fields, object, fileName)
}