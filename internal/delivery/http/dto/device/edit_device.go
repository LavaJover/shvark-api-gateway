package device

type EditDeviceRequest struct {
    EditDeviceParams EditDeviceParams `json:"editDeviceParams" binding:"required"`
}

type EditDeviceParams struct {
    DeviceName string `json:"deviceName"`
    Enabled    bool   `json:"enabled"`
}

type EditDeviceResponse struct {
    DeviceID   string `json:"deviceId,omitempty"`
    DeviceName string `json:"deviceName,omitempty"`
    Enabled    bool   `json:"enabled"`
}