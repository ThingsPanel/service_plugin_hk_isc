package services

// import (
// 	"encoding/json"
// 	"fmt"
// )

// // EncodeDeviceSearchRequest 定义请求结构，使用指针表示非必输字段
// type EncodeDeviceSearchRequest struct {
// 	Name             *string      `json:"name,omitempty"`
// 	RegionIndexCodes []string     `json:"regionIndexCodes,omitempty"`
// 	IsSubRegion      *bool        `json:"isSubRegion,omitempty"`
// 	PageNo           int          `json:"pageNo"`
// 	PageSize         int          `json:"pageSize"`
// 	AuthCodes        []string     `json:"authCodes,omitempty"`
// 	CapabilitySet    []string     `json:"capabilitySet,omitempty"`
// 	Expressions      []Expression `json:"expressions,omitempty"`
// 	OrderBy          *string      `json:"orderBy,omitempty"`
// 	OrderType        *string      `json:"orderType,omitempty"`
// }

// type Expression struct {
// 	Key      string   `json:"key"`
// 	Operator int      `json:"operator"`
// 	Values   []string `json:"values"`
// }

// // 辅助函数，用于创建字符串指针
// func StringPtr(s string) *string {
// 	return &s
// }

// // 辅助函数，用于创建布尔值指针
// func BoolPtr(b bool) *bool {
// 	return &b
// }

// // EncodeDeviceSearchResponse 定义响应结构
// type EncodeDeviceSearchResponse struct {
// 	Code string `json:"code"`
// 	Msg  string `json:"msg"`
// 	Data struct {
// 		Total    int            `json:"total"`
// 		PageNo   int            `json:"pageNo"`
// 		PageSize int            `json:"pageSize"`
// 		List     []EncodeDevice `json:"list"`
// 	} `json:"data"`
// }

// type EncodeDevice struct {
// 	BelongIndexCode string `json:"belongIndexCode"`
// 	Capability      string `json:"capability"`
// 	DeviceKey       string `json:"deviceKey"`
// 	DeviceType      string `json:"deviceType"`
// 	DevSerialNum    string `json:"devSerialNum"`
// 	DeviceCode      string `json:"deviceCode"`
// 	IndexCode       string `json:"indexCode"`
// 	Manufacturer    string `json:"manufacturer"`
// 	Name            string `json:"name"`
// 	RegionIndexCode string `json:"regionIndexCode"`
// 	RegionPath      string `json:"regionPath"`
// 	ResourceType    string `json:"resourceType"`
// 	TreatyType      string `json:"treatyType"`
// 	CreateTime      string `json:"createTime"`
// 	UpdateTime      string `json:"updateTime"`
// }

// func SearchEncodeDevices(request EncodeDeviceSearchRequest) (*EncodeDeviceSearchResponse, error) {
// 	const ARTEMIS_PATH = "/artemis"
// 	searchAPI := ARTEMIS_PATH + "/api/resource/v2/encodeDevice/search"

// 	body, err := json.Marshal(request)
// 	if err != nil {
// 		return nil, fmt.Errorf("error marshaling request: %v", err)
// 	}

// 	contentType := "application/json"
// 	responseBody, err := DoPostStringArtemis(searchAPI, body, contentType)
// 	if err != nil {
// 		return nil, fmt.Errorf("error making API request: %v", err)
// 	}

// 	var response EncodeDeviceSearchResponse
// 	err = json.Unmarshal([]byte(responseBody), &response)
// 	if err != nil {
// 		return nil, fmt.Errorf("error unmarshaling response: %v", err)
// 	}

// 	return &response, nil
// }

// // ExampleSearchEncodeDevices 展示如何使用修改后的结构
// func ExampleSearchEncodeDevices() {
// 	// request := EncodeDeviceSearchRequest{
// 	// 	Name:             StringPtr("test"),
// 	// 	RegionIndexCodes: []string{"8fyc8qw8280y0y43"},
// 	// 	IsSubRegion:      BoolPtr(true),
// 	// 	PageNo:           1,
// 	// 	PageSize:         10,
// 	// 	AuthCodes:        []string{"view"},
// 	// 	Expressions: []Expression{
// 	// 		{
// 	// 			Key:      "indexCode",
// 	// 			Operator: 0,
// 	// 			Values:   []string{"ayd8y80y1y082ye01y2e"},
// 	// 		},
// 	// 	},
// 	// 	OrderBy:   StringPtr("name"),
// 	// 	OrderType: StringPtr("desc"),
// 	// }

// 	request := EncodeDeviceSearchRequest{
// 		PageNo:   1,
// 		PageSize: 10,
// 		// Expressions: []Expression{
// 		// 	{
// 		// 		Key:      "indexCode",
// 		// 		Operator: 0,
// 		// 		Values:   []string{"51f18976a5704854b31d137e26beb136"},
// 		// 	},
// 		// },
// 	}

// 	response, err := SearchEncodeDevices(request)
// 	if err != nil {
// 		fmt.Printf("Error searching encode devices: %v\n", err)
// 		return
// 	}
// 	fmt.Println(response.Data)

// 	fmt.Printf("Total devices: %d\n", response.Data.Total)
// 	for _, device := range response.Data.List {
// 		fmt.Printf("Device Name: %s, Index Code: %s\n", device.Name, device.IndexCode)
// 	}
// }
