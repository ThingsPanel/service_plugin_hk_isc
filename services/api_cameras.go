package services

import (
	"encoding/json"
	"fmt"
	"service_hk_isc/model"

	"github.com/sirupsen/logrus"
)

/*
分页获取监控点资源

名称：分页获取监控点资源
描述:获取监控点列表接口可用来全量同步监控点信息，返回结果分页展示。
分组：视频资源
版本支持：V1.0
协议：HTTPS
请求路径：/api/resource/v1/cameras
URL：https://218.6.43.28:442/artemis/api/resource/v1/cameras
HTTP METHOD：POST
安全验证：API网关安全验证
*/

// CameraResourceRequest 定义请求结构
type CameraResourceRequest struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}

// CameraResourceResponse 定义响应结构
type CameraResourceResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Total    int      `json:"total"`
		PageNo   int      `json:"pageNo"`
		PageSize int      `json:"pageSize"`
		List     []Camera `json:"list"`
	} `json:"data"`
}

// Camera 定义监控点结构
type Camera struct {
	Altitude                  *string `json:"altitude"`
	CameraIndexCode           string  `json:"cameraIndexCode"`
	CameraName                string  `json:"cameraName"`
	CameraType                int     `json:"cameraType"`
	CameraTypeName            string  `json:"cameraTypeName"`
	CapabilitySet             string  `json:"capabilitySet"`
	CapabilitySetName         string  `json:"capabilitySetName"`
	IntelligentSet            string  `json:"intelligentSet"`
	IntelligentSetName        *string `json:"intelligentSetName"`
	ChannelNo                 string  `json:"channelNo"`
	ChannelType               string  `json:"channelType"`
	ChannelTypeName           *string `json:"channelTypeName"`
	CreateTime                *string `json:"createTime"`
	EncodeDevIndexCode        string  `json:"encodeDevIndexCode"`
	EncodeDevResourceType     *string `json:"encodeDevResourceType"`
	EncodeDevResourceTypeName *string `json:"encodeDevResourceTypeName"`
	GbIndexCode               *string `json:"gbIndexCode"`
	InstallLocation           string  `json:"installLocation"`
	KeyBoardCode              string  `json:"keyBoardCode"`
	Latitude                  *string `json:"latitude"`
	Longitude                 *string `json:"longitude"`
	Pixel                     int     `json:"pixel"`
	Ptz                       int     `json:"ptz"`
	PtzController             int     `json:"ptzController"`
	PtzControllerName         string  `json:"ptzControllerName"`
	PtzName                   string  `json:"ptzName"`
	RecordLocation            string  `json:"recordLocation"`
	RecordLocationName        string  `json:"recordLocationName"`
	RegionIndexCode           string  `json:"regionIndexCode"`
	Status                    int     `json:"status"`
	StatusName                string  `json:"statusName"`
	TransType                 int     `json:"transType"`
	TransTypeName             string  `json:"transTypeName"`
	TreatyType                string  `json:"treatyType"`
	TreatyTypeName            string  `json:"treatyTypeName"`
	Viewshed                  string  `json:"viewshed"`
	UpdateTime                string  `json:"updateTime"`
}

// GetCameraResources 实现分页获取监控点资源的 API 调用
func GetCameraResources(request CameraResourceRequest, config model.Voucher) (*CameraResourceResponse, error) {
	const ARTEMIS_PATH = "/artemis"
	apiPath := ARTEMIS_PATH + "/api/resource/v1/cameras"

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	contentType := "application/json"
	responseBody, err := DoPostStringArtemis(apiPath, body, contentType, config)
	if err != nil {
		return nil, fmt.Errorf("error making API request: %v", err)
	}

	var response CameraResourceResponse
	err = json.Unmarshal([]byte(responseBody), &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &response, nil
}

// 获取三方视频资源列表
func GetThirdCameraResources(request CameraResourceRequest, voucher model.Voucher) (any, error) {
	// 调用 GetCameraResources 函数获取监控点资源
	response, err := GetCameraResources(request, voucher)
	if err != nil {
		return nil, err
	}
	// 判断是否成功
	if response.Code != "0" {
		return nil, fmt.Errorf("API request failed: %s", response.Msg)
	}
	// 转换为三方视频资源列表
	var thirdDeviceList []model.ThirdDevice
	for _, camera := range response.Data.List {
		thirdDevice := model.ThirdDevice{
			DeviceName:   camera.CameraName,
			Description:  camera.CameraName,
			DeviceNumber: camera.CameraIndexCode,
		}
		thirdDeviceList = append(thirdDeviceList, thirdDevice)
	}

	// 构建三方视频资源列表响应
	rsp := model.ThirdDeviceRsp{
		List:  thirdDeviceList,
		Total: response.Data.Total,
	}
	return rsp, nil

}

// ExampleGetCameraResources 展示如何使用 GetCameraResources 函数
func ExampleGetCameraResources() {
	request := CameraResourceRequest{
		PageNo:   1,
		PageSize: 2,
	}
	config := model.Voucher{
		Host:      "218.6.43.28:442",
		AppKey:    "e82i01",
		AppSecret: "hdJivxLm2SpTw1qhkGL1",
	}

	response, err := GetCameraResources(request, config)
	if err != nil {
		fmt.Printf("Error getting camera resources: %v\n", err)
		return
	}
	logrus.Info("API Response: ", response)
	fmt.Printf("Total cameras: %d\n", response.Data.Total)
	for _, camera := range response.Data.List {
		fmt.Printf("Camera Name: %s, Index Code: %s\n", camera.CameraName, camera.CameraIndexCode)
	}
}
