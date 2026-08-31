package models

type DiscoveredCamera struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"` // ex: http://192.168.1.50/onvif/device_service
}
