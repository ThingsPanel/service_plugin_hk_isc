package httpservice

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"service_hk_isc/model"
	"service_hk_isc/services"
	"strconv"

	tpprotocolsdkgo "github.com/ThingsPanel/tp-protocol-sdk-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var HttpClient *tpprotocolsdkgo.Client

func Init() {
	go start()
}

func start() {
	var handler tpprotocolsdkgo.Handler = tpprotocolsdkgo.Handler{
		OnDisconnectDevice: OnDisconnectDevice,
		OnGetForm:          OnGetForm,
	}
	addr := viper.GetString("http_server.address")
	log.Println("http服务启动：", addr)
	err := handler.ListenAndServe(addr)
	if err != nil {
		log.Println("ListenAndServe() failed, err: ", err)
		return
	}
}

// OnGetForm 获取协议插件的json表单
func OnGetForm(w http.ResponseWriter, r *http.Request) {
	logrus.Info("OnGetForm")
	r.ParseForm() //解析参数，默认是不会解析的
	logrus.Info("【收到api请求】path", r.URL.Path)
	logrus.Info("query", r.URL.Query())

	//device_type := r.URL.Query()["device_type"][0]
	form_type := r.URL.Query()["form_type"][0]
	// service_identifier := r.URL.Query()["protocol_type"][0]
	// 根据需要对服务标识符进行验证，可不验证
	// if service_identifier != "xxxx" {
	// 	RspError(w, fmt.Errorf("not support protocol type: %s", service_identifier))
	// 	return
	// }
	//CFG配置表单 VCR凭证表单 SVCR服务凭证表单
	switch form_type {
	case "VCR":
		RspSuccess(w, nil)
	case "SVCR":
		//服务凭证类型表单
		RspSuccess(w, readFormConfigByPath("./form_service_voucher.json"))
	case "CFG":
		RspSuccess(w, nil)
	default:
		RspError(w, errors.New("not support form type: "+form_type))
	}
}

func OnDisconnectDevice(w http.ResponseWriter, r *http.Request) {
	logrus.Info("OnDisconnectDevice")
	r.ParseForm() //解析参数，默认是不会解析的
	logrus.Info("【收到api请求】path", r.URL.Path)
	logrus.Info("query", r.URL.Query())
	// 断开设备

	//RspSuccess(w, nil)
}

// ./form_config.json
func readFormConfigByPath(path string) interface{} {
	filePtr, err := os.Open(path)
	if err != nil {
		logrus.Warn("文件打开失败...", err.Error())
		return nil
	}
	defer filePtr.Close()
	var info interface{}
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&info)
	if err != nil {
		logrus.Warn("解码失败", err.Error())
		return info
	} else {
		logrus.Info("读取文件[form_config.json]成功...")
		return info
	}
}

func OnNotifyEvent(w http.ResponseWriter, r *http.Request) {
	logrus.Info("OnNotifyEvent")
	r.ParseForm() //解析参数，默认是不会解析的
	logrus.Info("【收到api请求】path", r.URL.Path)
	logrus.Info("query", r.Body)
	// 读取body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Warn("读取body失败", err.Error())
		return
	}
	logrus.Info("body", string(body))
	type NotifyEvent struct {
		MessageType string `json:"message_type"`
		Message     string `json:"message"`
	}
	// 解析到NotifyEvent
	var notifyEvent NotifyEvent
	err = json.Unmarshal(body, &notifyEvent)
	if err != nil {
		logrus.Warn("解析body失败", err.Error())
		RspError(w, err)
		return
	}
	logrus.Info("notifyEvent", notifyEvent)
	if notifyEvent.MessageType != "1" {
		type NotifyEventData struct {
			ServiceAccessID string `json:"service_access_id"`
		}
		var notifyEventData NotifyEventData
		err = json.Unmarshal([]byte(notifyEvent.Message), &notifyEventData)
		if err != nil {
			logrus.Warn("解析message失败", err.Error())
			RspError(w, err)
			return
		}
		return
	}
	RspSuccess(w, nil)
	// 处理事件通知
	//RspSuccess(w, nil)
}

// 处理设备列表请求
func OnGetDeviceList(w http.ResponseWriter, r *http.Request) {
	logrus.Info("OnGetDeviceList")

	// 使用r.URL.Query()一次性解析所有查询参数
	query := r.URL.Query()

	// 使用默认值和错误处理来获取查询参数
	voucherStr := query.Get("voucher")
	if voucherStr == "" {
		logrus.Warn("Missing voucher parameter")
		RspError(w, errors.New("missing voucher parameter"))

		return
	}

	// 转换页码和页面大小，使用默认值
	pageSize, err := strconv.Atoi(query.Get("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10 // 默认页面大小
	}

	pageNo, err := strconv.Atoi(query.Get("page_no"))
	if err != nil || pageNo <= 0 {
		pageNo = 1 // 默认页码
	}

	// 解析凭证
	var voucher model.Voucher
	if err := json.Unmarshal([]byte(voucherStr), &voucher); err != nil {
		logrus.Warn("Failed to parse voucher:", err)
		RspError(w, err)
		return
	}

	// 构建请求
	req := services.CameraResourceRequest{
		PageNo:   pageNo,
		PageSize: pageSize,
	}

	// 获取设备列表
	data, err := services.GetThirdCameraResources(req, voucher)
	if err != nil {
		logrus.Warn("Failed to get device list:", err)
		RspError(w, err)
		return
	}

	RspSuccess(w, data)
}
