package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	"github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/merchant"
	orderpb "github.com/LavaJover/shvark-order-service/proto/gen"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MerchantHandler struct {
	OrderClient *client.OrderClient
	WalletClient *client.HTTPWalletClient
}

func NewMerchanHandler(
	orderClient *client.OrderClient,
	walletClient *client.HTTPWalletClient,
) *MerchantHandler {
	return &MerchantHandler{
		OrderClient: orderClient,
		WalletClient: walletClient,
	}
}

// @Summary Create new deposit order
// @Description Create new pay-in order
// @Tags merchant
// @Accept json
// @Produce json
// @Param accountID path string true "merchant account ID"
// @Param input body merchant.CreatePayInRequest true "new deposit order details"
// @Success 201 {object} merchant.CreatePayInResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /merchant/order/{accountID}/deposit [post]
func (h *MerchantHandler) CreatePayIn(c *gin.Context) {
	merchantID := c.Param("accountID")
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "accountID path param missed"})
		return
	}
	var request merchant.CreatePayInRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orderServiceRequest := orderpb.CreateOrderRequest{
		MerchantId: merchantID,
		AmountFiat: request.Amount,
		Currency: request.Currency,
		Country: "Russia",
		ClientId: "",
		ExpiresAt: timestamppb.New(time.Now().Add(20*time.Minute)),
		MerchantOrderId: request.IternalID,
		Shuffle: 0,
		CallbackUrl: request.CallbackUrl,
		BankCode: request.Issuer,
		NspkCode: request.NspkCode,
	}
	if request.IsSbp {
		orderServiceRequest.PaymentSystem = "SBP"
	}else {
		orderServiceRequest.PaymentSystem = "C2C"
	}
	orderServiceResponse, err := h.OrderClient.CreateOrder(&orderServiceRequest)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, merchant.CreatePayInResponse{
		OrderID: orderServiceResponse.Order.OrderId,
		CardNumber: orderServiceResponse.Order.BankDetail.CardNumber,
		PhoneNumber: orderServiceResponse.Order.BankDetail.Phone,
		HolderName: orderServiceResponse.Order.BankDetail.Owner,
		Issuer: orderServiceResponse.Order.BankDetail.BankCode,
		NspkCode: orderServiceResponse.Order.BankDetail.NspkCode,
		Amount: orderServiceResponse.Order.AmountFiat,
		AmountByCurrency: orderServiceResponse.Order.AmountCrypto,
		CurrencyRate: orderServiceResponse.Order.CryptoRubRate,
		TimeExpires: orderServiceResponse.Order.ExpiresAt.Seconds,
	})
}

func (h *MerchantHandler) GetOrders(c *gin.Context) {
	
}

func (h *MerchantHandler) GetOrderStatus(c *gin.Context) {

}

// @Summary Get balance
// @Description Get account balance
// @Tags merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} merchant.GetAccountBalanceResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /merchant/accounts/balance [get]
func (h *MerchantHandler) GetAccountBalance(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
		return
	}

	userIDstr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID is not a string"})
		return
	}

	balance, err := h.WalletClient.GetBalance(userIDstr)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, merchant.GetAccountBalanceResponse{
		Balances: []merchant.AccountBalance{
			merchant.AccountBalance{
				Name: "test-name",
				Balance: strconv.FormatFloat(balance, 'f', 6, 64),
				Currency: "USDT",
			},
		},
	})

}

// @Summary Withdraw crypto from account wallet
// @Description Withdraw USDT
// @Tags merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body merchant.WithdrawRequest true "withdraw data"
// @Success 201 {object} merchant.WithdrawResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /merchant/accounts/withdraw/create [post]
func (h *MerchantHandler) Withdraw(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
		return
	}
	userIDstr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID is not a string"})
		return
	}
	var withdrawRequest merchant.WithdrawRequest
	if err := c.ShouldBindJSON(&withdrawRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txHash, err := h.WalletClient.Withdraw(
		userIDstr,
		withdrawRequest.ToAddress,
		withdrawRequest.Amount,
	)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, merchant.WithdrawResponse{
		TxHash: txHash,
		Amount: withdrawRequest.Amount,
		Currency: withdrawRequest.Currency,
		ToAddress: withdrawRequest.ToAddress,
	})

}