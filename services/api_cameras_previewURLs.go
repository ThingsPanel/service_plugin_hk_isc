package services

import (
	"encoding/json"
	"fmt"
	"service_hk_isc/model"

	"github.com/sirupsen/logrus"
)

/*
获取监控点预览取流URLv2

名称：获取监控点预览取流URLv2
描述:1.平台正常运行；平台已经添加过设备和监控点信息。 2.平台需要安装mgc取流服务。 3.三方平台通过openAPI获取到监控点数据，依据自身业务开发监控点导航界面。 4.调用本接口获取预览取流URL，协议类型包括：hik、rtsp、rtmp、hls。 5.通过开放平台的开发包进行实时预览或者使用标准的GUI播放工具进行实时预览。 6.为保证数据的安全性，取流URL设有有效时间，有效时间为5分钟。
分组：视频功能
版本支持：V1.4
请求基础定义
协议：HTTPS
请求路径：/api/video/v2/cameras/previewURLs
URL：https://218.6.43.28:442/artemis/api/video/v2/cameras/previewURLs
HTTP METHOD：POST
安全验证：API网关安全验证
*/

// CameraPreviewURLRequest 表示获取监控点预览取流URL的请求
type CameraPreviewURLRequest struct {
	CameraIndexCode string `json:"cameraIndexCode"`
	StreamType      int    `json:"streamType,omitempty"`
	Protocol        string `json:"protocol,omitempty"`
	Transmode       int    `json:"transmode,omitempty"`
	Expand          string `json:"expand,omitempty"`
	Streamform      string `json:"streamform,omitempty"`
}

// CameraPreviewURLResponse 表示获取监控点预览取流URL的响应
type CameraPreviewURLResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		URL string `json:"url"`
	} `json:"data"`
}

// GetCameraPreviewURL 实现获取监控点预览取流URL的 API 调用
func GetCameraPreviewURL(request CameraPreviewURLRequest, config model.Voucher) (*CameraPreviewURLResponse, error) {
	const ARTEMIS_PATH = "/artemis"
	apiPath := ARTEMIS_PATH + "/api/video/v2/cameras/previewURLs"

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	contentType := "application/json"
	responseBody, err := DoPostStringArtemis(apiPath, body, contentType, config)
	if err != nil {
		return nil, fmt.Errorf("error making API request: %v", err)
	}

	var response CameraPreviewURLResponse
	err = json.Unmarshal([]byte(responseBody), &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &response, nil
}

// GetCameraPreviewURLWrapper 获取监控点预览取流URL的包装函数
func GetCameraPreviewURLWrapper(request CameraPreviewURLRequest, voucher model.Voucher) (string, error) {
	response, err := GetCameraPreviewURL(request, voucher)
	if err != nil {
		return "", err
	}
	if response.Code != "0" {
		return "", fmt.Errorf("API request failed: %s", response.Msg)
	}
	return response.Data.URL, nil
}

// ExampleGetCameraPreviewURL 展示如何使用 GetCameraPreviewURL 函数
func ExampleGetCameraPreviewURL() {
	request := CameraPreviewURLRequest{
		CameraIndexCode: "0206ef351a9b42689c3e1de0e0827f5b",
		Protocol:        "rtsp",
		Expand:          "transcode=1&streamform=rtp",
	}
	config := model.Voucher{
		Host:      "127.0.0.1:442",
		AppKey:    "xxx",
		AppSecret: "xxxxxx",
	}

	response, err := GetCameraPreviewURL(request, config)
	if err != nil {
		fmt.Printf("Error getting camera preview URL: %v\n", err)
		return
	}
	logrus.Info("API Response: ", response)
	if response.Code == "0" {
		fmt.Printf("Preview URL: %s\n", response.Data.URL)
	} else {
		fmt.Printf("API request failed: %s - %s\n", response.Code, response.Msg)
	}
}
