package handlers

import (
	"net/http"
	"time"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	paymentRequest "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/payment/request"
	orderpb "github.com/LavaJover/shvark-order-service/proto/gen"
	"google.golang.org/protobuf/types/known/timestamppb"

	paymentResponse "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/payment/response"
	// "github.com/LavaJover/shvark-api-gateway/internal/domain"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	OrderClient *client.OrderClient
}

func NewPaymentHandler(orderClient *client.OrderClient) (*PaymentHandler, error) {
	return &PaymentHandler{
		OrderClient: orderClient,
	}, nil
}

// @Summary Create new H2H Pay-In
// @Description Create new Pay-In using host-to-host method
// @Tags payments
// @Accept json
// @Produce json
// @Param input body paymentRequest.CreateH2HPayInRequest true "pay-in info"
// @Success 200 {object} paymentResponse.CreateH2HPayInResponse
// @Failure 400 {object} paymentResponse.BadRequestErrorResponse
// @Failure 502 {object} paymentResponse.NoBankDetailsErrorResponse
// @Router /payments/in/h2h [post]
func (h *PaymentHandler) CreateH2HPayIn(c *gin.Context) {
	var payInRequest paymentRequest.CreateH2HPayInRequest
	if err := c.ShouldBindJSON(&payInRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ttl, err := time.ParseDuration(payInRequest.TTL)
	if err != nil {
		c.JSON(http.StatusBadRequest, paymentResponse.BadRequestErrorResponse{Error: err.Error()})
		return
	}
	response, err := h.OrderClient.CreateOrder(&orderpb.CreateOrderRequest{
		MerchantId: payInRequest.MerchantID,
		AmountFiat: payInRequest.AmountFiat,
		Currency: payInRequest.Currency,
		ClientId: payInRequest.ClientID,
		PaymentSystem: payInRequest.PaymentSystem,
		ExpiresAt: timestamppb.New(time.Now().Add(ttl)),
		MerchantOrderId: payInRequest.MerchantOrderID,
		Shuffle: payInRequest.Shuffle,
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, paymentResponse.NoBankDetailsErrorResponse{Error: "No available bank details"})
		return
	}

	c.JSON(http.StatusCreated, paymentResponse.CreateH2HPayInResponse{
		OrderID: response.Order.OrderId,
		AmountFiat: response.Order.AmountFiat,
		AmountCrypto: response.Order.AmountCrypto,
		Currency: response.Order.BankDetail.Currency,
		PaymentSystem: response.Order.BankDetail.PaymentSystem,
		Status: response.Order.Status,
		MerchantOrderID: response.Order.MerchantOrderId,
		CallbackURL: "http://callback/url",
		PaymentDetails: paymentResponse.PaymentDetails{
			CardNumber: response.Order.BankDetail.CardNumber,
			Owner: response.Order.BankDetail.Owner,
			Phone: response.Order.BankDetail.Phone,
			BankID: "bank-id-sbp-system",
			Bank: "user-friendly-bank",
			BankName: response.Order.BankDetail.BankName,
		},
		ExpiresAt: response.Order.ExpiresAt.String(),
	})
}

func (h *PaymentHandler) GetH2HPayInInfo(c *gin.Context) {

}

func (h *PaymentHandler) CancelPayIn(c *gin.Context) {

}

func (h *PaymentHandler) OpenPayInArbitrage(c *gin.Context) {

}

func (h *PaymentHandler) CreateRedirectPayIn(c *gin.Context) {

}

func (h *PaymentHandler) GetRedirectPayInInfo(c *gin.Context) {

}