package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	bankingpb "github.com/LavaJover/shvark-banking-service/proto/gen"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/durationpb"

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
// @Security BearerAuth
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

	bankDetailRequest := bankingpb.CreateBankDetailRequest{
		TraderId: request.TraderID,
		Country: request.Country,
		Currency: request.Currency,
		InflowCurrency: request.InflowCurrency,
		MinAmount: request.MinAmount,
		MaxAmount: request.MaxAmount,
		BankName: request.BankName,
		PaymentSystem: request.PaymentSystem,
		Delay: durationpb.New(delay),
		Enabled: request.Enabled,
		CardNumber: request.CardNumber,
		Phone: request.Phone,
		Owner: request.Owner,
		MaxOrdersSimultaneosly: request.MaxOrdersSimultaneosly,
		MaxAmountDay: request.MaxAmountDay,
		MaxAmountMonth: request.MaxAmountMonth,
		MaxQuantityDay: float64(request.MaxQuantityDay),
		MaxQuantityMonth: float64(request.MaxQuantityMonth),
		DeviceId: request.DeviceID,
	}

	response, err := h.BankingClient.CreateBankDetail(&bankDetailRequest)
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
// @Security BearerAuth
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
			MinAmount: response.BankDetail.MinAmount,
			MaxAmount: response.BankDetail.MaxAmount,
			BankName: response.BankDetail.BankName,
			PaymentSystem: response.BankDetail.PaymentSystem,
			Delay: response.BankDetail.Delay.AsDuration().Milliseconds(),
			Enabled: response.BankDetail.Enabled,
			CardNumber: response.BankDetail.CardNumber,
			Phone: response.BankDetail.Phone,
			Owner: response.BankDetail.Owner,
			MaxOrdersSimultaneosly: response.BankDetail.MaxOrdersSimultaneosly,
			MaxAmountDay: response.BankDetail.MaxAmountDay,
			MaxAmountMonth: response.BankDetail.MaxAmountMonth,
			MaxQuantityDay: int32(response.BankDetail.MaxQuantityDay),
			MaxQuantityMonth: int32(response.BankDetail.MaxQuantityMonth),
			DeviceID: response.BankDetail.DeviceId,
			InflowCurrency: response.BankDetail.InflowCurrency,
		},
	})
}

// @Summary Update bank detail
// @Description Update bank detail
// @Tags banking
// @Security BearerAuth
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

	bankDetailRequest := bankingpb.UpdateBankDetailRequest{
		BankDetail: &bankingpb.BankDetail{
			BankDetailId: request.BankDetail.ID,
			TraderId: request.BankDetail.TraderID,
			Currency: request.BankDetail.Currency,
			Country: request.BankDetail.Country,
			MinAmount: request.BankDetail.MinAmount,
			MaxAmount: request.BankDetail.MaxAmount,
			BankName: request.BankDetail.BankName,
			PaymentSystem: request.BankDetail.PaymentSystem,
			Enabled: request.BankDetail.Enabled,
			Delay: durationpb.New(delay),
			CardNumber: request.BankDetail.CardNumber,
			Phone: request.BankDetail.Phone,
			Owner: request.BankDetail.Owner,
			MaxOrdersSimultaneosly: request.BankDetail.MaxOrdersSimultaneosly,
			MaxAmountDay: request.BankDetail.MaxAmountDay,
			MaxAmountMonth: request.BankDetail.MaxAmountMonth,
			MaxQuantityDay: float64(request.BankDetail.MaxQuantityDay),
			MaxQuantityMonth: float64(request.BankDetail.MaxQuantityMonth),
			DeviceId: request.BankDetail.DeviceID,
			InflowCurrency: request.BankDetail.InflowCurrency,
		},
	}

	_, err = h.BankingClient.UpdateBankDetail(&bankDetailRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bankingResponse.UpdateBankDetailResponse{})
}

// @Summary Get bank details by trader ID
// @Description Get bank details by trader ID
// @Tags banking
// @Security BearerAuth
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
			MinAmount: bankDetail.MinAmount,
			MaxAmount: bankDetail.MaxAmount,
			BankName: bankDetail.BankName,
			PaymentSystem: bankDetail.PaymentSystem,
			Delay: bankDetail.Delay.AsDuration().Milliseconds(),
			Enabled: bankDetail.Enabled,
			CardNumber: bankDetail.CardNumber,
			Phone: bankDetail.Phone,
			Owner: bankDetail.Owner,
			MaxOrdersSimultaneosly: bankDetail.MaxOrdersSimultaneosly,
			MaxAmountDay: bankDetail.MaxAmountDay,
			MaxAmountMonth: bankDetail.MaxAmountMonth,
			MaxQuantityDay: int32(bankDetail.MaxQuantityDay),
			MaxQuantityMonth: int32(bankDetail.MaxQuantityMonth),
			DeviceID: bankDetail.DeviceId,
			InflowCurrency: bankDetail.InflowCurrency,
		}
	}

	c.JSON(http.StatusOK, bankingResponse.GetBankDetailsByTraderIDResponse{
		BankDetails: bankDetails,
	})
}

// @Summary Delete bank detail
// @Descrition Delete bank detail
// @Tags banking
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body bankingRequest.DeleteBankDetailRequest true "bank detail ID"
// @Success 200 {object} bankingResponse.DeleteBankDetailResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /banking/details/delete [post]
func (h *BankingHandler) DeleteBankDetail(c *gin.Context) {
	var request bankingRequest.DeleteBankDetailRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bankDetailID := request.BankDetailID
	response, err := h.BankingClient.DeleteBankDetail(bankDetailID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bankingResponse.DeleteBankDetailResponse{
		BankDetail: bankingResponse.BankDetail{
			ID: response.BankDetail.BankDetailId,
			TraderID: response.BankDetail.TraderId,
			Country: response.BankDetail.Country,
			Currency: response.BankDetail.Currency,
			MinAmount: response.BankDetail.MinAmount,
			MaxAmount: response.BankDetail.MaxAmount,
			BankName: response.BankDetail.BankName,
			PaymentSystem: response.BankDetail.PaymentSystem,
			Delay: response.BankDetail.Delay.AsDuration().Milliseconds(),
			Enabled: response.BankDetail.Enabled,
			CardNumber: response.BankDetail.CardNumber,
			Phone: response.BankDetail.Phone,
			Owner: response.BankDetail.Owner,
			MaxOrdersSimultaneosly: response.BankDetail.MaxOrdersSimultaneosly,
			MaxAmountDay: response.BankDetail.MaxAmountDay,
			MaxAmountMonth: response.BankDetail.MaxAmountMonth,
			MaxQuantityDay: int32(response.BankDetail.MaxQuantityDay),
			MaxQuantityMonth: int32(response.BankDetail.MaxQuantityMonth),
			DeviceID: response.BankDetail.DeviceId,
			InflowCurrency: response.BankDetail.InflowCurrency,
		},
	})
}