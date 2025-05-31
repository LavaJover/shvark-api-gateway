package handlers

import (
	"net/http"
	"time"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	"github.com/LavaJover/shvark-api-gateway/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
// @Param input body CreateBankDetailRequest true "New bank details"
// @Success 201 {object} CreateBankDetailResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /banking/details [post]
func (h *BankingHandler) CreateBankDetail(c *gin.Context) {
	var request CreateBankDetailRequest
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
	}

	response, err := h.BankingClient.CreateBankDetail(&bankDetail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateBankDetailResponse{
		BankDetailId: response.BankDetailId,
		Message: response.Message,
	})
}

type CreateBankDetailRequest struct {
	TraderID 		string			`json:"trader_id"`
	Country 		string			`json:"country"`
	Currency 		string			`json:"currency"`
	MinAmount 		float32			`json:"min_amount"`
	MaxAmount 		float32			`json:"max_amount"`
	BankName 		string			`json:"bank_name"`
	PaymentSystem 	string			`json:"payment_system"`
	Delay			string			`json:"delay" example:"100s"`
	Enabled 		bool			`json:"enabled"`
}

type CreateBankDetailResponse struct {
	BankDetailId string `json:"bank_detail_id"`
	Message 	 string `json:"message"`
}

// @Summary Get bank detail by ID
// @Description Get bank detail by ID
// @Tags banking
// @Accept json
// @Produce json
// @Param uuid path string true "bank detail UUID"
// @Success 200 {object} GetBankDetailByIDResponse
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

	c.JSON(http.StatusOK, GetBankDetailByIDResponse{
		BankDetail: BankDetail{
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
		},
	})
}

type BankDetail struct {
	ID 				string	`json:"id"`
	TraderID 		string	`json:"trader_id"`
	Country 		string	`json:"country"`
	Currency 		string	`json:"currency"`
	MinAmount 		float32	`json:"min_amount"`
	MaxAmount 		float32	`json:"max_amount"`
	BankName 		string	`json:"bank_name"`
	PaymentSystem 	string	`json:"payment_system"`
	Delay			string	`json:"delay"`
	Enabled 		bool	`json:"enabled"`
}

type GetBankDetailByIDResponse struct {
	BankDetail	`json:"bank_detail"`
}

func (h *BankingHandler) UpdateBankDetail(c *gin.Context) {
	
}