package model

type ThirdDevice struct {
	DeviceName   string `json:"device_name"`
	Description  string `json:"description"`
	DeviceNumber string `json:"device_number"`
}

type ThirdDeviceRsp struct {
	List  []ThirdDevice `json:"list"`
	Total int           `json:"total"`
}
