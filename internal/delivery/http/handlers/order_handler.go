package handlers

import (
	"github.com/LavaJover/shvark-api-gateway/internal/client"
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

func (h *OrderHandler) CreateOrder(c *gin.Context) {

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