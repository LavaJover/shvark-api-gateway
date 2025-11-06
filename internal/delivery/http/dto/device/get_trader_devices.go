package device

type GetTraderDevicesRequest struct {
    TraderID string `json:"traderId"`
}

type GetTraderDevicesResponse struct {
    Devices []Device `json:"devices"`
}

type Device struct {
    DeviceID   string `json:"deviceId"`
    DeviceName string `json:"deviceName"`
    TraderID   string `json:"traderId"`
    Enabled    bool   `json:"enabled"`
}