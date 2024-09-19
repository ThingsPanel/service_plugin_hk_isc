package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"io"
	"net/http"
	"service_hk_isc/model"
	"time"
)

// ArtemisConfig 存储配置信息
type ArtemisConfig struct {
	Host      string
	AppKey    string
	AppSecret string
}

// 全局配置
// var config = ArtemisConfig{
// 	Host:      "218.6.43.28:442",
// 	AppKey:    "e82i01",
// 	AppSecret: "hdJivxLm2SpTw1qhkGL1",
// }

// 生成签名
func generateSignature(secret, method, accept, contentType, date, path string) string {
	stringToSign := method + "\n" + accept + "\n" + contentType + "\n" + date + "\n" + path
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// 发送 POST 请求
func DoPostStringArtemis(path string, body []byte, contentType string, config model.Voucher) (string, error) {
	url := "https://" + config.Host + path
	method := "POST"
	accept := "application/json"
	date := time.Now().UTC().Format(http.TimeFormat)

	signature := generateSignature(config.AppSecret, method, accept, contentType, date, path)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", accept)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Date", date)
	req.Header.Set("X-Ca-Key", config.AppKey)
	req.Header.Set("X-Ca-Signature", signature)
	req.Header.Set("X-Ca-Signature-Headers", "x-ca-key")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}
