package device

type EditDeviceRequest struct {
	EditDeviceParams EditDeviceParams `json:"editDeviceParams"` 
}

type EditDeviceParams struct {
	DeviceName string `json:"deviceName"`
	Enabled    bool   `json:"enabled"`
}

type EditDeviceResponse struct {

}