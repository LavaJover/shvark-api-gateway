package handlers

import (
	"net/http"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	"github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/device"
	orderpb "github.com/LavaJover/shvark-order-service/proto/gen"
	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	OrderClient *client.OrderClient
}

func NewDeviceHandler(orderClient *client.OrderClient) (*DeviceHandler, error) {
	return &DeviceHandler{
		OrderClient: orderClient,
	}, nil
}

// @Summary Create new device
// @Description Create new logic device for atomatic
// @Tags devices
// @Accept json
// @Produce json
// @Param input body device.CreateDeviceRequest true "new device data"
// @Success 201 {object} device.CreateDeviceResponse
// @Failure 500 {object} ErrorResponse
// @Router /devices [post] 
func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	var request device.CreateDeviceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := h.OrderClient.CreateDevice(
		&orderpb.CreateDeviceRequest{
			DeviceName: request.DeviceName,
			TraderId: request.TraderID,
			Enabled: request.Enabled,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create new device"})
		return
	}

	c.JSON(http.StatusCreated, device.CreateDeviceResponse{})
}

// @Summary Get trader devices
// @Description Get trader devices
// @Tags devices
// @Accept json
// @Produce json
// @Param traderId path string true "trader ID"
// @Success 200 {object} device.GetTraderDevicesResponse
// @Failure 500 {object} ErrorResponse
// @Router /devices/{traderId} [get]
func (h *DeviceHandler) GetTraderDevices(c *gin.Context) {
	traderID := c.Param("traderId")
	resp, err := h.OrderClient.GetTraderDevices(
		&orderpb.GetTraderDevicesRequest{
			TraderId: traderID,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	devices := make([]device.Device, len(resp.Devices))
	for i, d := range resp.Devices {
		devices[i] = device.Device{
			DeviceID: d.DeviceId,
			DeviceName: d.DeviceName,
			TraderID: d.TraderId,
			Enabled: d.Enabled,
		}
	}

	c.JSON(http.StatusOK, device.GetTraderDevicesResponse{
		Devices: devices,
	})
}

// @Summary Delete device
// @Description Delete exact device by id
// @Tags devices
// @Accept json
// @Produce json
// @Param deviceId path string true "device ID"
// @Success 200 {object} device.DeleteDeviceResponse
// @Failure 500 {object} ErrorResponse
// @Router /devices/{deviceId} [delete]
func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	deviceID := c.Param("deviceId")
	_, err := h.OrderClient.DeleteDevice(
		&orderpb.DeleteDeviceRequest{
			DeviceId: deviceID,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, device.DeleteDeviceResponse{})
}

// @Summary Edit device
// @Description Edit device params
// @Tags devices
// @Accept json
// @Produce json
// @Param input body device.EditDeviceRequest true "device edit params"
// @Param deviceId path string true "ID of device to edit"
// @Success 200 {object} device.EditDeviceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /devices/{deviceId}/edit [patch]
func (h *DeviceHandler) EditDevice(c *gin.Context) {
	deviceID := c.Param("deviceId")
	var request device.EditDeviceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.OrderClient.EditeDevice(
		&orderpb.EditDeviceRequest{
			DeviceId: deviceID,
			Params: &orderpb.EditDeviceParams{
				DeviceName: request.EditDeviceParams.DeviceName,
				Enabled: request.EditDeviceParams.Enabled,
			},
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, device.EditDeviceResponse{})
}