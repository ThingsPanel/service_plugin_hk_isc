package services

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
