package handlers

import (
	"net/http"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	orderRequest "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/order/request"
	orderResponse "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/order/response"
	"github.com/LavaJover/shvark-api-gateway/internal/domain"
	"github.com/gin-gonic/gin"
)

var _ = orderRequest.CreateOrderRequest{}

type OrderHandler struct {
	OrderClient *client.OrderClient
}

func NewOrderHandler(addr string) (*OrderHandler, error) {
	orderClient, err := client.NewOrderClient(addr)
	if err != nil {
		return nil, err
	}

	return &OrderHandler{
		OrderClient: orderClient,
	}, nil
}

// @Summary Create new Pay-In order
// @Description Create new Pay-In order
// @Tags orders
// @Accept json
// @Produce json
// @Param input body orderRequest.CreateOrderRequest true "new order details"
// @Success 200 {object} orderResponse.CreateOrderResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var request orderRequest.CreateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderRequest := domain.Order{
		MerchantID: request.MerchantID,
		AmountFiat: request.AmountFiat,
		Currency: request.Currency,
		Country: request.Country,
		ClientEmail: request.ClientData,
		MetadataJSON: request.Metadata,
		PaymentSystem: request.PaymentSystem,
		ExpiresAt: request.ExpiresAt,
	}

	response, err := h.OrderClient.CreateOrder(&orderRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orderResponse.CreateOrderResponse{
		OrderID: response.Order.OrderId,
		OrderStatus: response.Order.Status,
		AmountFiat: response.Order.AmountFiat,
		AmountCrypto: response.Order.AmountCrypto,
		BankDetail: orderResponse.BankDetail{
			ID: response.Order.BankDetail.BankDetailId,
			TraderID: response.Order.BankDetail.TraderId,
			Currency: response.Order.BankDetail.Currency,
			Country: response.Order.BankDetail.Country,
			MinAmount: response.Order.BankDetail.MinAmount,
			MaxAmount: response.Order.BankDetail.MaxAmount,
			BankName: response.Order.BankDetail.BankName,
			PaymentSystem: response.Order.BankDetail.PaymentSystem,
			Enabled: response.Order.BankDetail.Enabled,
			Delay: response.Order.BankDetail.Delay.String(),
		},
	})
}