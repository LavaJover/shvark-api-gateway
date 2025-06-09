package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	"github.com/LavaJover/shvark-api-gateway/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	bankingRequest "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/banking/request"
	bankingResponse "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/banking/response"
)

type BankingHandler struct {
	BankingClient *client.BankingClient
}

func NewBankingHandler(addr string) (*BankingHandler, error) {
	bankingClient, err := client.NewBankingClient(addr)
	if err != nil {
		return nil, err
	}

	return &BankingHandler{
		BankingClient: bankingClient,
	}, nil
}

// @Summary Create new bank detail
// @Description Create new bank detail
// @Tags banking
// @Accept json
// @Produce json
// @Param input body bankingRequest.CreateBankDetailRequest true "New bank details"
// @Success 201 {object} bankingResponse.CreateBankDetailResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /banking/details [post]
func (h *BankingHandler) CreateBankDetail(c *gin.Context) {
	var request bankingRequest.CreateBankDetailRequest
	fmt.Println(c.Request.Body)
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delay, err := time.ParseDuration(request.Delay)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bankDetail := domain.BankDetail{
		TraderID: request.TraderID,
		Country: request.Country,
		Currency: request.Currency,
		MinAmount: request.MinAmount,
		MaxAmount: request.MaxAmount,
		BankName: request.BankName,
		PaymentSystem: request.PaymentSystem,
		Delay: delay,
		Enabled: request.Enabled,
		CardNumber: request.CardNumber,
		Phone: request.Phone,
		Owner: request.Owner,
		MaxOrdersSimultaneosly: request.MaxOrdersSimultaneosly,
		MaxAmountDay: request.MaxAmountDay,
		MaxAmountMonth: request.MaxAmountMonth,
		MaxQuantityDay: request.MaxQuantityDay,
		MaxQuantityMonth: request.MaxQuantityMonth,
		DeviceID: request.DeviceID,
	}

	response, err := h.BankingClient.CreateBankDetail(&bankDetail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bankingResponse.CreateBankDetailResponse{
		BankDetailId: response.BankDetailId,
		Message: response.Message,
	})
}

// @Summary Get bank detail by ID
// @Description Get bank detail by ID
// @Tags banking
// @Accept json
// @Produce json
// @Param uuid path string true "bank detail UUID"
// @Success 200 {object} bankingResponse.GetBankDetailByIDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /banking/details/{uuid} [get]
func (h *BankingHandler) GetBankDetailByID(c *gin.Context) {

	bankDetailID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.BankingClient.GetBankDetailByID(bankDetailID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bankingResponse.GetBankDetailByIDResponse{
		BankDetail: bankingResponse.BankDetail{
			ID: response.BankDetail.BankDetailId,
			TraderID: response.BankDetail.TraderId,
			Country: response.BankDetail.Country,
			Currency: response.BankDetail.Currency,
			MinAmount: float32(response.BankDetail.MinAmount),
			MaxAmount: float32(response.BankDetail.MaxAmount),
			BankName: response.BankDetail.BankName,
			PaymentSystem: response.BankDetail.PaymentSystem,
			Delay: response.BankDetail.Delay.String(),
			Enabled: response.BankDetail.Enabled,
			CardNumber: response.BankDetail.CardNumber,
			Phone: response.BankDetail.Phone,
			Owner: response.BankDetail.Owner,
			MaxOrdersSimultaneosly: response.BankDetail.MaxOrdersSimultaneosly,
			MaxAmountDay: int32(response.BankDetail.MaxAmountDay),
			MaxAmountMonth: int32(response.BankDetail.MaxAmountMonth),
			MaxQuantityDay: int32(response.BankDetail.MaxQuantityDay),
			MaxQuantityMonth: int32(response.BankDetail.MaxQuantityMonth),
			DeviceID: response.BankDetail.DeviceId,
		},
	})
}

// @Summary Update bank detail
// @Description Update bank detail
// @Tags banking
// @Accept json
// @Produce json
// @Param input body bankingRequest.UpdateBankDetailRequest true "New data for bank detail with given ID"
// @Success 200 {object} bankingResponse.UpdateBankDetailResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /banking/details [patch]
func (h *BankingHandler) UpdateBankDetail(c *gin.Context) {
	var request bankingRequest.UpdateBankDetailRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delay, err := time.ParseDuration(request.BankDetail.Delay)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedBankDetail := domain.BankDetail{
		ID: request.BankDetail.ID,
		TraderID: request.BankDetail.TraderID,
		Country: request.BankDetail.Country,
		Currency: request.BankDetail.Currency,
		MinAmount: request.BankDetail.MinAmount,
		MaxAmount: request.BankDetail.MaxAmount,
		BankName: request.BankDetail.BankName,
		PaymentSystem: request.BankDetail.PaymentSystem,
		Delay: delay,
		Enabled: request.BankDetail.Enabled,
		CardNumber: request.BankDetail.CardNumber,
		Phone: request.BankDetail.Phone,
		Owner: request.BankDetail.Owner,
		MaxOrdersSimultaneosly: request.BankDetail.MaxOrdersSimultaneosly,
		MaxAmountDay: request.BankDetail.MaxAmountDay,
		MaxAmountMonth: request.BankDetail.MaxAmountMonth,
		MaxQuantityDay: request.BankDetail.MaxQuantityDay,
		MaxQuantityMonth: request.BankDetail.MaxAmountMonth,
		DeviceID: request.BankDetail.DeviceID,
	}

	_, err = h.BankingClient.UpdateBankDetail(&updatedBankDetail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bankingResponse.UpdateBankDetailResponse{})
}

// @Summary Get bank details by trader ID
// @Description Get bank details by trader ID
// @Tags banking
// @Accept json
// @Produce json
// @Param trader query string false "trader uuid"
// @Success 200 {object} bankingResponse.GetBankDetailsByTraderIDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /banking/details/ [get]
func (h *BankingHandler) GetBankDetailsByTraderID(c *gin.Context) {
	traderID, err := uuid.Parse(c.Query("trader"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.BankingClient.GetBankDetailsByTraderID(traderID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bankDetails := make([]bankingResponse.BankDetail, len(response.BankDetails))
	for i, bankDetail := range response.BankDetails {
		bankDetails[i] = bankingResponse.BankDetail{
			ID: bankDetail.BankDetailId,
			TraderID: bankDetail.TraderId,
			Country: bankDetail.Country,
			Currency: bankDetail.Currency,
			MinAmount: float32(bankDetail.MinAmount),
			MaxAmount: float32(bankDetail.MaxAmount),
			BankName: bankDetail.BankName,
			PaymentSystem: bankDetail.PaymentSystem,
			Delay: bankDetail.Delay.String(),
			Enabled: bankDetail.Enabled,
			CardNumber: bankDetail.CardNumber,
			Phone: bankDetail.Phone,
			Owner: bankDetail.Owner,
			MaxOrdersSimultaneosly: bankDetail.MaxOrdersSimultaneosly,
			MaxAmountDay: int32(bankDetail.MaxAmountDay),
			MaxAmountMonth: int32(bankDetail.MaxAmountMonth),
			MaxQuantityDay: int32(bankDetail.MaxQuantityDay),
			MaxQuantityMonth: int32(bankDetail.MaxQuantityMonth),
			DeviceID: bankDetail.DeviceId,
		}
	}

	c.JSON(http.StatusOK, bankingResponse.GetBankDetailsByTraderIDResponse{
		BankDetails: bankDetails,
	})
}