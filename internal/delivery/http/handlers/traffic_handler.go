package handlers

import (
	"net/http"
	"strconv"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	orderpb "github.com/LavaJover/shvark-order-service/proto/gen"
	"github.com/gin-gonic/gin"
)

type TrafficHandler struct {
	orderClient *client.OrderClient
}

func NewTrafficHandler(
	orderClient *client.OrderClient,
) *TrafficHandler {
	return &TrafficHandler{
		orderClient: orderClient,
	}
}

// @Summary Set trader traffic lock status
// @Description On/Off trader traffic lock status
// @Tags traffic
// @Accept json
// @Produce json
// @Param traderID path string true "trader ID"
// @Param unlocked query bool true "is unlocked"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /traffic/traders/{traderID} [patch]
func (h *TrafficHandler) SetTraderLockTrafficStatus(c *gin.Context) {
	traderID := c.Param("traderID")
	if traderID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "traderID is required"})
		return
	}

	unlockedStr := c.Query("unlocked")
	unlocked, err := strconv.ParseBool(unlockedStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "unlocked must be a boolean value"})
		return
	}

	_, err = h.orderClient.SetTraderLockTrafficStatus(&orderpb.SetTraderLockTrafficStatusRequest{
		TraderId: traderID,
		Unlocked: unlocked,
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

// @Summary Set merchant traffic lock status
// @Description On/Off merchant traffic lock status
// @Tags traffic
// @Accept json
// @Produce json
// @Param merchantID path string true "merchant ID"
// @Param unlocked query bool true "is unlocked"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /traffic/merchants/{merchantID} [patch]
func (h *TrafficHandler) SetMerchantLockTrafficStatus(c *gin.Context) {
	merchantID := c.Param("merchantID")
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "merchantID is required"})
		return
	}

	unlockedStr := c.Query("unlocked")
	unlocked, err := strconv.ParseBool(unlockedStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "unlocked must be a boolean value"})
		return
	}

	_, err = h.orderClient.SetMerchantLockTrafficStatus(&orderpb.SetMerchantLockTrafficStatusRequest{
		MerchantId: merchantID,
		Ubnlocked: unlocked,
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

// @Summary Set manual traffic lock status
// @Description On/Off manual traffic lock status for specific traffic record
// @Tags traffic
// @Accept json
// @Produce json
// @Param trafficID path string true "traffic record ID"
// @Param unlocked query bool true "is unlocked"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /traffic/{trafficID}/manual [patch]
func (h *TrafficHandler) SetManuallyLockTrafficStatus(c *gin.Context) {
	trafficID := c.Param("trafficID")
	if trafficID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "trafficID is required"})
		return
	}

	unlockedStr := c.Query("unlocked")
	unlocked, err := strconv.ParseBool(unlockedStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "unlocked must be a boolean value"})
		return
	}

	_, err = h.orderClient.SetManuallyLockTrafficStatus(&orderpb.SetManuallyLockTrafficStatusRequest{
		TrafficId: trafficID,
		Unlocked:  unlocked,
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

// @Summary Set antifraud traffic lock status
// @Description On/Off antifraud traffic lock status
// @Tags traffic
// @Accept json
// @Produce json
// @Param traderID path string true "trader ID"
// @Param unlocked query bool true "is unlocked"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /traffic/antifraud/{traderID} [patch]
func (h *TrafficHandler) SetAntifraudLockTrafficStatus(c *gin.Context) {
	traderID := c.Param("traderID")
	if traderID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "traderID is required"})
		return
	}

	unlockedStr := c.Query("unlocked")
	unlocked, err := strconv.ParseBool(unlockedStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "unlocked must be a boolean value"})
		return
	}

	_, err = h.orderClient.SetAntifraudLockTrafficStatus(&orderpb.SetAntifraudLockTrafficStatusRequest{
		TraderId: traderID,
		Unlocked: unlocked,
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

// @Summary Get traffic lock statuses
// @Description Get all lock statuses for specific traffic record
// @Tags traffic
// @Accept json
// @Produce json
// @Param trafficID path string true "traffic record ID"
// @Success 200 {object} LockStatusesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /traffic/{trafficID}/lock-statuses [get]
func (h *TrafficHandler) GetTrafficLockStatuses(c *gin.Context) {
	trafficID := c.Param("trafficID")
	if trafficID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "trafficID is required"})
		return
	}

	response, err := h.orderClient.GetTrafficLockStatuses(&orderpb.GetTrafficLockStatusesRequest{
		TrafficId: trafficID,
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, LockStatusesResponse{
		TrafficID:         response.TrafficId,
		MerchantUnlocked:  response.MerchantUnlocked,
		TraderUnlocked:    response.TraderUnlocked,
		AntifraudUnlocked: response.AntifraudUnlocked,
		ManuallyUnlocked:  response.ManuallyUnlocked,
	})
}

// @Summary Check if traffic is unlocked
// @Description Check if traffic is unlocked by any method
// @Tags traffic
// @Accept json
// @Produce json
// @Param trafficID path string true "traffic record ID"
// @Success 200 {object} TrafficUnlockedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /traffic/{trafficID}/unlocked [get]
func (h *TrafficHandler) CheckTrafficUnlocked(c *gin.Context) {
	trafficID := c.Param("trafficID")
	if trafficID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "trafficID is required"})
		return
	}

	response, err := h.orderClient.CheckTrafficUnlocked(&orderpb.CheckTrafficUnlockedRequest{
		TrafficId: trafficID,
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, TrafficUnlockedResponse{
		TrafficID: response.TrafficId,
		Unlocked:  response.Unlocked,
	})
}

// Структуры ответов для HTTP
type LockStatusesResponse struct {
	TrafficID         string `json:"traffic_id"`
	MerchantUnlocked  bool   `json:"merchant_unlocked"`
	TraderUnlocked    bool   `json:"trader_unlocked"`
	AntifraudUnlocked bool   `json:"antifraud_unlocked"`
	ManuallyUnlocked  bool   `json:"manually_unlocked"`
}

type TrafficUnlockedResponse struct {
	TrafficID string `json:"traffic_id"`
	Unlocked  bool   `json:"unlocked"`
}