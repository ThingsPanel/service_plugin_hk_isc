package services

import (
	"context"
	"encoding/json"
	httpclient "service_hk_isc/http_client"
	"service_hk_isc/model"
	"service_hk_isc/mqtt"
	"sync"
	"time"

	"github.com/ThingsPanel/tp-protocol-sdk-go/api"
	"github.com/sirupsen/logrus"
)

var (
	ServiceAccessPointMap = make(map[string]*ServiceAccessPoint)
	mapMutex              sync.RWMutex
)

// 缓存服务接入点的信息
type ServiceAccessPoint struct {
	Data        *api.ServiceAccess // 服务访问数据
	ctx         context.Context    // 用于控制协程生命周期的上下文
	cancel      context.CancelFunc // 用于取消上下文的函数
	wg          sync.WaitGroup     // 用于等待协程完成的等待组
	stopChan    chan struct{}      // 用于发送停止信号的通道
	isRunning   bool               // 标记服务是否正在运行
	runningLock sync.Mutex         // 用于保护 isRunning 的互斥锁
}

func GetServiceAccessPointMap() (map[string]*ServiceAccessPoint, error) {
	logrus.Info("获取服务接入点列表")
	rspData, err := httpclient.GetServiceAccessList()
	if err != nil {
		return nil, err
	}

	mapMutex.Lock()
	defer mapMutex.Unlock()

	for _, v := range rspData.Data {
		if _, exists := ServiceAccessPointMap[v.ServiceAccessID]; !exists {
			ctx, cancel := context.WithCancel(context.Background())
			sap := &ServiceAccessPoint{
				Data:     &v,
				ctx:      ctx,
				cancel:   cancel,
				stopChan: make(chan struct{}),
			}
			ServiceAccessPointMap[v.ServiceAccessID] = sap
			if err := sap.Start(); err != nil {
				logrus.Errorf("Failed to start ServiceAccessPoint %s: %v", v.ServiceAccessID, err)
			}
		}
	}

	// Remove any ServiceAccessPoints that are no longer in the response
	for id := range ServiceAccessPointMap {
		if !containsServiceAccess(rspData.Data, id) {
			ServiceAccessPointMap[id].Stop()
			delete(ServiceAccessPointMap, id)
		}
	}

	return ServiceAccessPointMap, nil
}

func containsServiceAccess(data []api.ServiceAccess, id string) bool {
	for _, v := range data {
		if v.ServiceAccessID == id {
			return true
		}
	}
	return false
}

func (sap *ServiceAccessPoint) Start() error {
	sap.runningLock.Lock()
	defer sap.runningLock.Unlock()

	if sap.isRunning {
		return nil
	}

	sap.wg.Add(1)
	go func() {
		defer sap.wg.Done()
		sap.run()
	}()

	sap.isRunning = true
	return nil
}

func (sap *ServiceAccessPoint) Stop() {
	sap.runningLock.Lock()
	defer sap.runningLock.Unlock()

	if !sap.isRunning {
		return
	}

	sap.cancel()
	close(sap.stopChan)
	sap.wg.Wait()
	sap.isRunning = false
}

func (sap *ServiceAccessPoint) run() {
	logrus.Info("ServiceAccessPoint run start")

	sap.DoSomething() // 首次执行

	ticker := time.NewTicker(4 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-sap.ctx.Done():
			return
		case <-sap.stopChan:
			return
		case <-ticker.C:
			logrus.Info("ServiceAccessPoint run ...")
			sap.DoSomething()
		}
	}
}

// DoSomething 执行具体的操作
func (sap *ServiceAccessPoint) DoSomething() {
	for _, v := range sap.Data.Devices {
		deviceNumber := v.DeviceNumber
		logrus.Debug("设备编号：", deviceNumber)
		// 在这里添加对设备的具体操作
		request := CameraPreviewURLRequest{
			CameraIndexCode: v.DeviceNumber,
			Protocol:        "rtsp",
			Expand:          "transcode=1&streamform=rtp",
		}
		// 解析凭证
		voucher, err := ParseVoucher(v.Voucher)
		if err != nil {
			logrus.Errorf("Error parsing voucher: %v", err)
			continue
		}
		response, err := GetCameraPreviewURL(request, voucher)
		if err != nil {
			logrus.Errorf("Error getting camera preview URL: %v", err)
			continue
		}
		// 处理响应数据
		if response.Code != "0" {
			logrus.Errorf("Error getting camera preview URL: %s", response.Msg)
			continue
		}
		// 发送属性
		var payloadMap = make(map[string]interface{})
		if response.Data.URL == "" {
			logrus.Errorf("Error getting camera preview URL: %s", response.Msg)
			continue
		}
		payloadMap["previewUrl"] = response.Data.URL
		err = mqtt.PublishAttributes(v.ID, payloadMap)
		if err != nil {
			logrus.Errorf("Error publishing attributes: %v", err)
			continue
		}

	}
}

// 解析凭证信息
func ParseVoucher(voucher string) (model.Voucher, error) {
	var parsedVoucher model.Voucher
	err := json.Unmarshal([]byte(voucher), &parsedVoucher)
	if err != nil {
		return model.Voucher{}, err
	}
	return parsedVoucher, nil
}
