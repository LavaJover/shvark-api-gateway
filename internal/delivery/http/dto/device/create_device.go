package device

type CreateDeviceRequest struct {
	DeviceName 	string `json:"deviceName"`
	TraderID 	string `json:"traderId"`
	Enabled		bool   `json:"enabled"`
}

type CreateDeviceResponse struct {
	
}