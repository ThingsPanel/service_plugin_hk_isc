package services

type HkIsc struct {
}

func NewHkIsc() *HkIsc {
	return &HkIsc{}
}

func (h *HkIsc) Run() error {
	_, err := GetServiceAccessPointMap()
	if err != nil {
		return err
	}
	//ExampleGetCameraResources()

	//ExampleSearchEncodeDevices()
	return nil
}

func (h *HkIsc) Close() {
}
