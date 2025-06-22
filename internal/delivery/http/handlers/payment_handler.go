package handlers

import (
	"net/http"
	"time"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	paymentRequest "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/payment/request"
	orderpb "github.com/LavaJover/shvark-order-service/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
// @Success 201 {object} paymentResponse.CreateH2HPayInResponse
// @Failure 400 {object} paymentResponse.BadRequestErrorResponse
// @Failure 404 {object} paymentResponse.NoBankDetailsErrorResponse
// @Failure 409 {object} paymentResponse.ErrorResponse
// @Failure 502 {object} paymentResponse.ErrorResponse
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
		CallbackUrl: payInRequest.CallbackURL,
	})
	if err != nil  {
		if status, ok := status.FromError(err); ok {
			if status.Code() == codes.NotFound {
				c.JSON(http.StatusNotFound, paymentResponse.NoBankDetailsErrorResponse{Error: err.Error()})
				return
			}else if status.Code() == codes.FailedPrecondition {
				c.JSON(http.StatusConflict, paymentResponse.ErrorResponse{Error: err.Error()})
				return
			}else {
				c.JSON(http.StatusBadGateway, paymentResponse.ErrorResponse{Error: err.Error()})
			}
		}
	}

	c.JSON(http.StatusCreated, paymentResponse.CreateH2HPayInResponse{
		OrderID: response.Order.OrderId,
		AmountFiat: response.Order.AmountFiat,
		AmountCrypto: response.Order.AmountCrypto,
		Currency: response.Order.BankDetail.Currency,
		PaymentSystem: response.Order.BankDetail.PaymentSystem,
		Status: response.Order.Status,
		MerchantOrderID: response.Order.MerchantOrderId,
		CallbackURL: response.Order.CallbackUrl,
		TPayLink: "tpay/link",
		Recalculated: response.Order.Recalculated,
		CryptoRubRate: response.Order.CryptoRubRate,
		PaymentDetails: paymentResponse.PaymentDetails{
			CardNumber: response.Order.BankDetail.CardNumber,
			Owner: response.Order.BankDetail.Owner,
			Phone: response.Order.BankDetail.Phone,
			BankID: "bank-id-sbp-system",
			Bank: "user-friendly-bank",
			BankName: response.Order.BankDetail.BankName,
		},
		ExpiresAt: response.Order.ExpiresAt.Seconds,
	})
}

// @Summary Get H2h Pay-in info
// @Description Get host-to-host pay-in order info
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "order id"
// @Success 200 {object} paymentResponse.GetH2HPayInInfoResponse
// @Failure 400 {object} paymentResponse.BadRequestErrorResponse
// @Failure 404 {object} paymentResponse.ErrorResponse
// @Router /payments/in/h2h/{id} [get]
func (h *PaymentHandler) GetH2HPayInInfo(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, paymentResponse.BadRequestErrorResponse{Error: "id path param missed"})
		return
	}

	response, err := h.OrderClient.GetOrderByID(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, paymentResponse.ErrorResponse{Error: "Order info is unavailable now"})
		return
	}

	c.JSON(http.StatusOK, paymentResponse.GetH2HPayInInfoResponse{
		OrderID: response.Order.OrderId,
		AmountFiat: response.Order.AmountFiat,
		AmountCrypto: response.Order.AmountCrypto,
		Currency: response.Order.BankDetail.Currency,
		PaymentSystem: response.Order.BankDetail.PaymentSystem,
		Status: response.Order.Status,
		MerchantOrderID: response.Order.MerchantOrderId,
		CallbackURL: response.Order.CallbackUrl,
		Recalculated: response.Order.Recalculated,
		CryptoRubRate: response.Order.CryptoRubRate,
		PaymentDetails: paymentResponse.PaymentDetails{
			CardNumber: response.Order.BankDetail.CardNumber,
			Owner: response.Order.BankDetail.Owner,
			Phone: response.Order.BankDetail.Phone,
			BankID: "bank id in sbp",
			Bank: "user-friendly",
			BankName: response.Order.BankDetail.BankName,
		},
		ExpiresAt: response.Order.ExpiresAt.Seconds,
		TPayLink: "tpay-link",
	})

}

// @Summary Cancel Pay In order
// @Description Cancel Pay in order
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "order ID"
// @Success 200 {object} paymentResponse.CancelPayInResponse
// @Failure 400 {object} paymentResponse.ErrorResponse
// @Failure 502 {object} paymentResponse.ErrorResponse
// @Router /payments/in/h2h/{id}/cancel [post]
func (h *PaymentHandler) CancelPayIn(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, paymentResponse.ErrorResponse{Error: "order id path param missed"})
		return
	}

	response, err := h.OrderClient.CancelOrder(orderID)
	if err != nil {
		c.JSON(http.StatusBadGateway, paymentResponse.ErrorResponse{Error: "order with given ID was not found"})
		return
	}

	c.JSON(http.StatusOK, paymentResponse.CancelPayInResponse{
		Message: response.Message,
	})
}


func (h *PaymentHandler) OpenPayInArbitrage(c *gin.Context) {
	var request paymentRequest.CreateDisputeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func (h *PaymentHandler) CreateRedirectPayIn(c *gin.Context) {

}

func (h *PaymentHandler) GetRedirectPayInInfo(c *gin.Context) {

}