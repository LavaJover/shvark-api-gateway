package handlers

import (
	"net/http"
	"time"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	orderRequest "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/order/request"
	orderResponse "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/order/response"
	orderpb "github.com/LavaJover/shvark-order-service/proto/gen"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
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
// @Security BearerAuth
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

	ttl, err := time.ParseDuration(request.TTL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderRequest := orderpb.CreateOrderRequest{
		Shuffle: request.Shuffle,
		MerchantId: request.MerchantID,
		ClientId: request.ClientID,
		MerchantOrderId: request.MerchantOrderID,
		AmountFiat: request.AmountFiat,
		Currency: request.Currency,
		Country: request.Country,
		PaymentSystem: request.PaymentSystem,
		ExpiresAt: timestamppb.New(time.Now().Add(ttl)),
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

// @Summary Get order by orderID
// @Description Get order by orderID
// @Tags orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param uuid path string true "order UUID"
// @Success 200 {object} orderResponse.GetOrderByIDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /orders/{uuid} [get]
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.OrderClient.GetOrderByID(orderID.String())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orderResponse.GetOrderByIDResponse{
		Order: orderResponse.Order{
			OrderID: response.Order.OrderId,
			Status: response.Order.Status,
			AmountFiat: response.Order.AmountFiat,
			AmountCrypto: response.Order.AmountCrypto,
			ExpiresAt: response.Order.ExpiresAt.AsTime(),
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
				Owner: response.Order.BankDetail.Owner,
				CardNumber: response.Order.BankDetail.CardNumber,
				Phone: response.Order.BankDetail.Phone,
			},
		},
	})
}

// @Summary Get orders by trader ID
// @Description Get orders by trader ID
// @Tags orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param traderUUID path string true "Trader UUID path param"
// @Success 200 {object} orderResponse.GetOrdersByTraderIDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /orders/trader/{traderUUID} [get]
func (h *OrderHandler) GetOrdersByTraderID(c *gin.Context) {
	traderID, err := uuid.Parse(c.Param("traderUUID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.OrderClient.GetOrdersByTraderID(traderID.String())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	responseOrders := response.Orders
	orders := make([]orderResponse.Order, len(responseOrders))

	for i, responseOrder := range responseOrders {
		orders[i] = orderResponse.Order{
			OrderID: responseOrder.OrderId,
			Status: responseOrder.Status,
			AmountFiat: responseOrder.AmountFiat,
			AmountCrypto: responseOrder.AmountCrypto,
			ExpiresAt: responseOrder.ExpiresAt.AsTime(),
			BankDetail: orderResponse.BankDetail{
				ID: responseOrder.BankDetail.BankDetailId,
				TraderID: responseOrder.BankDetail.TraderId,
				Currency: responseOrder.BankDetail.Currency,
				Country: responseOrder.BankDetail.Country,
				MinAmount: responseOrder.BankDetail.MinAmount,
				MaxAmount: responseOrder.BankDetail.MaxAmount,
				BankName: responseOrder.BankDetail.BankName,
				PaymentSystem: responseOrder.BankDetail.PaymentSystem,
				Enabled: responseOrder.BankDetail.Enabled,
				Delay: responseOrder.BankDetail.Delay.String(),
				Owner: responseOrder.BankDetail.Owner,
				CardNumber: responseOrder.BankDetail.CardNumber,
				Phone: responseOrder.BankDetail.Phone,
			},
		}
	}

	c.JSON(http.StatusOK, orderResponse.GetOrdersByTraderIDResponse{
		Orders: orders,
	})
}

// @Summary Approve order by order uuid
// @Description Approve order by order uuid
// @Tags orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body orderRequest.ApproveOrderRequest true "Order UUID"
// @Success 200 {object} orderResponse.ApproveOrderResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /orders/approve [post]
func (h *OrderHandler) ApproveOrder(c *gin.Context) {
	var request orderRequest.ApproveOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orderID := request.OrderID
	response, err := h.OrderClient.ApproveOrder(orderID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orderResponse.ApproveOrderResponse{
		Message: response.Message,
	})
}

// @Summary Cancel order by order uuid
// @Description Cancel order by order uuid
// @Tags orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body orderRequest.CancelOrderRequest true "Order UUID"
// @Success 200 {object} orderResponse.CancelOrderResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /orders/cancel [post]
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	var request orderRequest.CancelOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orderID := request.OrderID
	response, err := h.OrderClient.CancelOrder(orderID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orderResponse.CancelOrderResponse{
		Message: response.Message,
	})
}

// @Summary Open order dispute by order uuid
// @Description Open order dispute by order uuid
// @Tags orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body orderRequest.OpenOrderDisputeRequest true "Order UUID"
// @Success 200 {object} orderResponse.OpenOrderDisputeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /orders/disputes/open [post]
func (h *OrderHandler) OpenOrderDispute(c *gin.Context) {
	var request orderRequest.OpenOrderDisputeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orderID := request.OrderID
	response, err := h.OrderClient.OpenOrderDispute(orderID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orderResponse.OpenOrderDisputeResponse{
		Message: response.Message,
	})
}

// @Summary Resolve order dispute by order uuid
// @Description Resolve order dispute by order uuid
// @Tags orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body orderRequest.ResolveOrderDisputeRequest true "Order UUID"
// @Success 200 {object} orderResponse.ResolveOrderDisputeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /orders/disputes/resolve [post]
func (h *OrderHandler) ResolveOrderDispute(c *gin.Context) {
	var request orderRequest.ResolveOrderDisputeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orderID := request.OrderID
	response, err := h.OrderClient.ResolveOrderDispute(orderID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orderResponse.ResolveOrderDisputeResponse{
		Message: response.Message,
	})
}