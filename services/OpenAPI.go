package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"service_hk_isc/model"

	"github.com/google/uuid"
)

func DoPostStringArtemis(path string, body []byte, contentType string, config model.Voucher) (string, error) {
	urlStr := fmt.Sprintf("https://%s%s", config.Host, path)
	method := "POST"
	accept := "application/json"
	date := time.Now().UTC().Format(http.TimeFormat)
	contentMD5 := calculateContentMD5(body)
	xCaTimestamp := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	xCaNonce := generateUUID()

	headers := map[string]string{
		"Accept":         accept,
		"Content-MD5":    contentMD5,
		"Content-Type":   contentType,
		"Date":           date,
		"x-ca-key":       config.AppKey,
		"x-ca-timestamp": xCaTimestamp,
		"x-ca-nonce":     xCaNonce,
	}

	log.Printf("Request Headers: %+v", headers)

	signature := generateSignature(config.AppSecret, method, headers, path, body)

	req, err := http.NewRequest(method, urlStr, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("X-Ca-Signature", signature)
	req.Header.Set("X-Ca-Signature-Headers", getSignatureHeaders(headers))

	log.Printf("Final Request Headers: %+v", req.Header)

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

	log.Printf("Response Body: %s", string(respBody))

	return string(respBody), nil
}

func generateSignature(appSecret, method string, headers map[string]string, path string, body []byte) string {
	stringToSign := buildStringToSign(method, headers, path, body)
	log.Printf("StringToSign:\n%s", stringToSign)

	h := hmac.New(sha256.New, []byte(appSecret))
	h.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	log.Printf("Generated Signature: %s", signature)

	return signature
}

func buildStringToSign(method string, headers map[string]string, path string, body []byte) string {
	accept := headers["Accept"]
	contentMD5 := headers["Content-MD5"]
	contentType := headers["Content-Type"]
	date := headers["Date"]

	headerString := buildHeaderString(headers)
	urlString := buildUrlString(path, body, headers["Content-Type"])

	parts := []string{
		method,
		accept,
		contentMD5,
		contentType,
		date,
	}

	// Remove empty parts
	var nonEmptyParts []string
	for _, part := range parts {
		if part != "" {
			nonEmptyParts = append(nonEmptyParts, part)
		}
	}

	// Add header string
	if headerString != "" {
		nonEmptyParts = append(nonEmptyParts, headerString)
	}

	// Add URL string
	nonEmptyParts = append(nonEmptyParts, urlString)

	return strings.Join(nonEmptyParts, "\n")
}

func buildHeaderString(headers map[string]string) string {
	var keys []string
	for k := range headers {
		lowerK := strings.ToLower(k)
		if strings.HasPrefix(lowerK, "x-ca-") && lowerK != "x-ca-signature" {
			keys = append(keys, lowerK)
		}
	}
	sort.Strings(keys)

	var headerParts []string
	for _, k := range keys {
		headerParts = append(headerParts, fmt.Sprintf("%s:%s", k, strings.TrimSpace(headers[k])))
	}

	if len(headerParts) > 0 {
		return strings.Join(headerParts, "\n")
	}
	return ""
}

func buildUrlString(path string, body []byte, contentType string) string {
	u, _ := url.Parse(path)
	query := u.Query()

	if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		bodyForm, _ := url.ParseQuery(string(body))
		for k, v := range bodyForm {
			if len(v) > 0 {
				query.Add(k, v[0])
			}
		}
	}

	var keys []string
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var queryParts []string
	for _, k := range keys {
		v := query.Get(k)
		if v != "" {
			queryParts = append(queryParts, fmt.Sprintf("%s=%s", k, v))
		} else {
			queryParts = append(queryParts, k)
		}
	}
	queryString := strings.Join(queryParts, "&")

	if queryString != "" {
		return fmt.Sprintf("%s?%s", u.Path, queryString)
	}
	return u.Path
}

func getSignatureHeaders(headers map[string]string) string {
	var keys []string
	for k := range headers {
		lowerK := strings.ToLower(k)
		if strings.HasPrefix(lowerK, "x-ca-") && lowerK != "x-ca-signature" {
			keys = append(keys, lowerK)
		}
	}
	sort.Strings(keys)
	return strings.Join(keys, ",")
}

func calculateContentMD5(body []byte) string {
	hash := md5.Sum(body)
	return base64.StdEncoding.EncodeToString(hash[:])
}

func generateUUID() string {
	return uuid.New().String()
}
