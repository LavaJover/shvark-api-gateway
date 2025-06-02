package handlers

import (
	"net/http"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	"github.com/LavaJover/shvark-api-gateway/internal/domain"
	"github.com/gin-gonic/gin"
)

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
// @Param input body CreateOrderRequest true "new order details"
// @Success 200 {object} CreateOrderResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var request CreateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderRequest := domain.Order{
		MerchantID: request.MerchantID,
		Amount: request.Amount,
		Currency: request.Currency,
		Country: request.Country,
		ClientEmail: request.ClientEmail,
		MetadataJSON: request.Metadata,
		PaymentSystem: request.PaymentSystem,
	}

	response, err := h.OrderClient.CreateOrder(&orderRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, OrderResponse{
		OrderID: response.Order.OrderId,
		OrderStatus: response.Order.Status.String(),
	})
}

type CreateOrderRequest struct {
	MerchantID  	string   `json:"merchant_id"`
	Amount	    	float32	 `json:"amount"`
	Currency   		string   `json:"currency"`
	Country	   		string 	 `json:"country"`
	ClientEmail 	string	 `json:"client_email"`
	Metadata    	string 	 `json:"metadata"`
	PaymentSystem   string   `json:"payment_system"`
}

type OrderResponse struct {
	OrderID     string `json:"order_id"`
	OrderStatus string `json:"order_status"`
}

type CreateOrderResponse struct {

}