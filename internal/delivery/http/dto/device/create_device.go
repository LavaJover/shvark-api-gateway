package device

type CreateDeviceRequest struct {
    DeviceName  string `json:"deviceName" binding:"required"`
    TraderID    string `json:"traderId" binding:"required"`
    Enabled     bool   `json:"enabled"`
}

type CreateDeviceResponse struct {
    DeviceID   string `json:"deviceId,omitempty"`
    DeviceName string `json:"deviceName,omitempty"`
    TraderID   string `json:"traderId,omitempty"`
    Enabled    bool   `json:"enabled"`
}