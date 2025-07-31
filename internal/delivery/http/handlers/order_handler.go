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
		Shuffle:         request.Shuffle,
		MerchantId:      request.MerchantID,
		ClientId:        request.ClientID,
		MerchantOrderId: request.MerchantOrderID,
		AmountFiat:      request.AmountFiat,
		Currency:        request.Currency,
		Country:         request.Country,
		PaymentSystem:   request.PaymentSystem,
		ExpiresAt:       timestamppb.New(time.Now().Add(ttl)),
	}

	response, err := h.OrderClient.CreateOrder(&orderRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orderResponse.CreateOrderResponse{
		OrderID:      response.Order.OrderId,
		OrderStatus:  response.Order.Status,
		AmountFiat:   response.Order.AmountFiat,
		AmountCrypto: response.Order.AmountCrypto,
		BankDetail: orderResponse.BankDetail{
			ID:            response.Order.BankDetail.BankDetailId,
			TraderID:      response.Order.BankDetail.TraderId,
			Currency:      response.Order.BankDetail.Currency,
			Country:       response.Order.BankDetail.Country,
			MinAmount:     response.Order.BankDetail.MinAmount,
			MaxAmount:     response.Order.BankDetail.MaxAmount,
			BankName:      response.Order.BankDetail.BankName,
			PaymentSystem: response.Order.BankDetail.PaymentSystem,
			Enabled:       response.Order.BankDetail.Enabled,
			Delay:         response.Order.BankDetail.Delay.String(),
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
			OrderID:      response.Order.OrderId,
			Status:       response.Order.Status,
			AmountFiat:   response.Order.AmountFiat,
			AmountCrypto: response.Order.AmountCrypto,
			ExpiresAt:    response.Order.ExpiresAt.AsTime(),
			TraderReward: response.Order.TraderRewardPercent,
			CreatedAt:    response.Order.CreatedAt.AsTime(),
			UpdatedAt:    response.Order.UpdatedAt.AsTime(),
			CryptoRubRate:   response.Order.CryptoRubRate,
			MerchantOrderID: response.Order.MerchantOrderId,
			MerchantID: response.Order.MerchantId,
			BankDetail: orderResponse.BankDetail{
				ID:            response.Order.BankDetail.BankDetailId,
				TraderID:      response.Order.BankDetail.TraderId,
				Currency:      response.Order.BankDetail.Currency,
				Country:       response.Order.BankDetail.Country,
				MinAmount:     response.Order.BankDetail.MinAmount,
				MaxAmount:     response.Order.BankDetail.MaxAmount,
				BankName:      response.Order.BankDetail.BankName,
				PaymentSystem: response.Order.BankDetail.PaymentSystem,
				Enabled:       response.Order.BankDetail.Enabled,
				Delay:         response.Order.BankDetail.Delay.String(),
				Owner:         response.Order.BankDetail.Owner,
				CardNumber:    response.Order.BankDetail.CardNumber,
				Phone:         response.Order.BankDetail.Phone,
			},
		},
	})
}

// @Summary Get order by merchant order ID
// @Description Get order by merchant order ID
// @Tags orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "order id in merchant system"
// @Success 200 {object} orderResponse.GetOrderByIDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /orders/merchant/{id} [get]
func (h *OrderHandler) GetOrderByMerchantOrderID(c *gin.Context) {
	merchantOrderID := c.Param("id")
	if merchantOrderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id path param missed"})
		return
	}
	response, err := h.OrderClient.GetOrderByMerchantOrderID(merchantOrderID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orderResponse.GetOrderByIDResponse{
		Order: orderResponse.Order{
			OrderID:      response.Order.OrderId,
			Status:       response.Order.Status,
			AmountFiat:   response.Order.AmountFiat,
			AmountCrypto: response.Order.AmountCrypto,
			ExpiresAt:    response.Order.ExpiresAt.AsTime(),
			TraderReward: response.Order.TraderRewardPercent,
			CreatedAt:    response.Order.CreatedAt.AsTime(),
			UpdatedAt:    response.Order.UpdatedAt.AsTime(),
			CryptoRubRate:   response.Order.CryptoRubRate,
			MerchantOrderID: response.Order.MerchantOrderId,
			MerchantID: response.Order.MerchantId,
			BankDetail: orderResponse.BankDetail{
				ID:            response.Order.BankDetail.BankDetailId,
				TraderID:      response.Order.BankDetail.TraderId,
				Currency:      response.Order.BankDetail.Currency,
				Country:       response.Order.BankDetail.Country,
				MinAmount:     response.Order.BankDetail.MinAmount,
				MaxAmount:     response.Order.BankDetail.MaxAmount,
				BankName:      response.Order.BankDetail.BankName,
				PaymentSystem: response.Order.BankDetail.PaymentSystem,
				Enabled:       response.Order.BankDetail.Enabled,
				Delay:         response.Order.BankDetail.Delay.String(),
				Owner:         response.Order.BankDetail.Owner,
				CardNumber:    response.Order.BankDetail.CardNumber,
				Phone:         response.Order.BankDetail.Phone,
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
// @Param page query int false "Page number" default(1) minimum(1)
// @Param limit query int false "Items per page" default(10) minimum(1) maximum(100)
// @Param sort_by query string false "Sort field (amount_fiat, expires_at)" default(expires_at)
// @Param sort_order query string false "Sort order (asc, desc)" default(desc)
// @Param status query []string false "Filter by status" collectionFormat(multi)
// @Param min_amount query number false "Minimum amount"
// @Param max_amount query number false "Maximum amount"
// @Param date_from query string false "Date from (YYYY-MM-DD)"
// @Param date_to query string false "Date to (YYYY-MM-DD)"
// @Param currency query string false "Currency code"
// @Param order_id query string false "Order ID"
// @Param merchant_order_id query string false "Merchant order ID"
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
	var request orderRequest.OrderQueryParams
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dateFrom, dateTo time.Time
	if request.DateFrom != "" {
		dateFrom, _ = time.Parse("2006-01-02", request.DateFrom)
	}
	if request.DateTo != "" {
		dateTo, _ = time.Parse("2006-01-02", request.DateTo)
	}

	response, err := h.OrderClient.GetOrdersByTraderID(
		&orderpb.GetOrdersByTraderIDRequest{
			TraderId:  traderID.String(),
			Page:      request.Page,
			Limit:     request.Limit,
			SortBy:    request.SortBy,
			SortOrder: request.SortOrder,
			Filters: &orderpb.OrderFilters{
				Statuses:      request.Status,
				MinAmountFiat: request.MinAmount,
				MaxAmountFiat: request.MaxAmount,
				DateFrom:      timestamppb.New(dateFrom),
				DateTo:        timestamppb.New(dateTo),
				Currency:      request.Currency,
			},
		},
	)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	responseOrders := response.Orders
	orders := make([]orderResponse.Order, len(responseOrders))

	for i, responseOrder := range responseOrders {
		orders[i] = orderResponse.Order{
			OrderID:      responseOrder.OrderId,
			Status:       responseOrder.Status,
			AmountFiat:   responseOrder.AmountFiat,
			AmountCrypto: responseOrder.AmountCrypto,
			ExpiresAt:    responseOrder.ExpiresAt.AsTime(),
			TraderReward: responseOrder.TraderRewardPercent,
			CreatedAt:    responseOrder.CreatedAt.AsTime(),
			UpdatedAt:    responseOrder.UpdatedAt.AsTime(),
			CryptoRubRate:   responseOrder.CryptoRubRate,
			MerchantOrderID: responseOrder.MerchantOrderId,
			MerchantID: responseOrder.MerchantId,
			BankDetail: orderResponse.BankDetail{
				ID:            responseOrder.BankDetail.BankDetailId,
				TraderID:      responseOrder.BankDetail.TraderId,
				Currency:      responseOrder.BankDetail.Currency,
				Country:       responseOrder.BankDetail.Country,
				MinAmount:     responseOrder.BankDetail.MinAmount,
				MaxAmount:     responseOrder.BankDetail.MaxAmount,
				BankName:      responseOrder.BankDetail.BankName,
				PaymentSystem: responseOrder.BankDetail.PaymentSystem,
				Enabled:       responseOrder.BankDetail.Enabled,
				Delay:         responseOrder.BankDetail.Delay.String(),
				Owner:         responseOrder.BankDetail.Owner,
				CardNumber:    responseOrder.BankDetail.CardNumber,
				Phone:         responseOrder.BankDetail.Phone,
			},
		}
	}

	c.JSON(http.StatusOK, orderResponse.GetOrdersByTraderIDResponse{
		Orders: orders,
		Pagination: orderResponse.Pagination{
			CurrentPage:  response.Pagination.CurrentPage,
			TotalPages:   response.Pagination.TotalPages,
			TotalItems:   response.Pagination.TotalItems,
			ItemsPerPage: response.Pagination.ItemsPerPage,
		},
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

// @Summary 	Get order statistics
// @Description Get order statistics
// @Tags 		orders
// @Security 	BearerAuth
// @Accept 		json
// @Produce 	json
// @Param       date_from query string true "Дата начала (RFC3339 format, e.g. 2025-07-21T00:00:00Z)"
// @Param       date_to   query string true "Дата конца (RFC3339 format, e.g. 2025-07-21T23:59:59Z)"
// @Success 200 {object} orderResponse.GetOrderStatsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /orders/statistics [get]
func (h *OrderHandler) GetOrderStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
		return
	}
	userIDstr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userID"})
		return
	}
	dateFromStr := c.Query("date_from")
	dateToStr := c.Query("date_to")

	dateFrom, err := time.Parse(time.RFC3339, dateFromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date_from format, expected RFC3339"})
		return
	}

	dateTo, err := time.Parse(time.RFC3339, dateToStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date_to format, expected RFC3339"})
		return
	}
	resp, err := h.OrderClient.GetOrderStats(
		userIDstr,
		dateFrom,
		dateTo,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no stats was found"})
		return
	}
	c.JSON(http.StatusOK, orderResponse.GetOrderStatsResponse{
		TotalOrders: resp.TotalOrders,
		SucceedOrders: resp.SucceedOrders,
		CanceledOrders: resp.CanceledOrders,
		ProcessedAmountFiat: float64(resp.ProcessedAmountFiat),
		ProcessedAmountCrypto: float64(resp.ProcessedAmountCrypto),
		CanceledAmountFiat: float64(resp.CanceledAmountFiat),
		CanceledAmountCrypto: float64(resp.CanceledAmountCrypto),
		IncomeCrypto: float64(resp.IncomeCrypto),
	})
}