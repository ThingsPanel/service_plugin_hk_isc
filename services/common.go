package services

import (
	"fmt"
	httpclient "service_hk_isc/http_client"

	"github.com/ThingsPanel/tp-protocol-sdk-go/api"
	"github.com/sirupsen/logrus"
)

// 认证设备并获取设备信息
func AuthDevice(deviceSecret string) (deviceInfo *api.DeviceConfigResponse, err error) {
	voucher := AssembleVoucher(deviceSecret)
	// 读取设备信息
	deviceInfo, err = httpclient.GetDeviceConfig(voucher, "")
	if err != nil {
		// 获取设备信息失败，请检查连接包是否正确
		logrus.Error(err)
		return
	}
	if deviceInfo.Code != 200 {
		err = fmt.Errorf("device auth failed, code: %d, message: %s", deviceInfo.Code, deviceInfo.Message)
		logrus.Error(err)
	}
	return
}

// 凭证信息组装
func AssembleVoucher(deviceSecret string) (voucher string) {
	return fmt.Sprintf(`{"UID":"%s"}`, deviceSecret)
}
