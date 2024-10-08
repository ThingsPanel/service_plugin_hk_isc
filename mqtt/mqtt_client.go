package mqtt

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	tpprotocolsdkgo "github.com/ThingsPanel/tp-protocol-sdk-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var MqttClient *tpprotocolsdkgo.MQTTClient

func InitClient() {
	log.Println("创建mqtt客户端")
	// 创建新的MQTT客户端实例
	addr := viper.GetString("mqtt.broker")
	username := viper.GetString("mqtt.username")
	password := viper.GetString("mqtt.password")
	client := tpprotocolsdkgo.NewMQTTClient(addr, username, password)
	// 尝试连接到MQTT代理
	if err := client.Connect(); err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	log.Println("连接成功")
	MqttClient = client
}

type MqttPayload struct {
	DeviceID string `json:"device_id"`
	Values   []byte `json:"values"`
}

// 获取消息id
func GetMessageID() string {
	// 获取当前Unix时间戳
	timestamp := time.Now().Unix()
	// 将时间戳转换为字符串
	timestampStr := strconv.FormatInt(timestamp, 10)
	// 截取后七位
	messageID := timestampStr[len(timestampStr)-7:]

	return messageID
}

// 组装payload{"device_id":device_id,"values":{key:value...}}
// values是base64编码的数据
func AssemblePayload(deviceID string, payload []byte) ([]byte, error) {
	var mqttPayload MqttPayload
	mqttPayload.DeviceID = deviceID
	mqttPayload.Values = payload
	newMsgJson, err := json.Marshal(mqttPayload)
	if err != nil {
		return nil, err
	}
	return newMsgJson, nil
}

// 发布遥测消息
func PublishTelemetry(deviceID string, data map[string]interface{}) error {
	topic := viper.GetString("mqtt.telemetry_topic_to_publish")
	qos := viper.GetUint("mqtt.qos")
	// map转json
	payload, err := json.Marshal(data)
	if err != nil {
		logrus.Warn("map转json失败:", err)
		return err
	}
	// 组装payload
	newMsgJson, err := AssemblePayload(deviceID, payload)
	if err != nil {
		logrus.Warn("组装payload失败:", err)
		return err
	}
	err = MqttClient.Publish(topic, string(newMsgJson), uint8(qos))
	if err != nil {
		logrus.Warn("发送消息失败:", err)
		return err
	}
	logrus.Debug("遥测主题:", topic)
	logrus.Debug("消息内容:", string(payload))
	logrus.Debug("\n==>tp 发送消息成功:", string(newMsgJson))

	return nil
}

// 发布命令响应
func PublishCommandResponse(deviceID string, messageID string, data map[string]interface{}) error {
	topic := viper.GetString("mqtt.command_response_topic_to_publish")
	qos := viper.GetUint("mqtt.qos")
	// map转json
	payload, err := json.Marshal(data)
	if err != nil {
		logrus.Warn("map转json失败:", err)
		return err
	}
	// 组装payload
	newMsgJson, err := AssemblePayload(deviceID, payload)
	if err != nil {
		logrus.Warn("组装payload失败:", err)
		return err
	}
	// 组装主题
	topic = topic + messageID
	err = MqttClient.Publish(topic, string(newMsgJson), uint8(qos))
	if err != nil {
		logrus.Warn("发送消息失败:", err)
		return err
	}
	logrus.Debug("命令响应主题:", topic)
	logrus.Debug("消息内容:", string(payload))
	logrus.Debug("\n==>tp 发送消息成功:", string(newMsgJson))

	return nil
}

// 发布属性消息
func PublishAttributes(deviceID string, data map[string]interface{}) error {
	topic := viper.GetString("mqtt.attributes_topic_to_publish") + GetMessageID()
	qos := viper.GetUint("mqtt.qos")
	// map转json
	payload, err := json.Marshal(data)
	if err != nil {
		logrus.Warn("map转json失败:", err)
		return err
	}
	// 组装payload
	newMsgJson, err := AssemblePayload(deviceID, payload)
	if err != nil {
		logrus.Warn("组装payload失败:", err)
		return err
	}
	err = MqttClient.Publish(topic, string(newMsgJson), uint8(qos))
	if err != nil {
		logrus.Warn("发送消息失败:", err)
		return err
	}
	logrus.Debug("属性主题:", topic)
	logrus.Debug("消息内容:", string(payload))
	logrus.Debug("\n==>tp 发送消息成功:", string(newMsgJson))

	return nil
}
