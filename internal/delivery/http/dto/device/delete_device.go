package device

type DeleteDeviceRequest struct {
    DeviceID string `json:"deviceId"`
}

type DeleteDeviceResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message,omitempty"`
}