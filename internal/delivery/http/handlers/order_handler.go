package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	"github.com/LavaJover/shvark-api-gateway/internal/common"
	orderRequest "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/order/request"
	orderResponse "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/order/response"
	orderpb "github.com/LavaJover/shvark-order-service/proto/gen/order"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc/status"
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
				OrderId: 	   request.OrderID,
				MerchantOrderId: request.MerchantOrderID,
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
			DeviceId: responseOrder.BankDetail.DeviceId,
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

// GetAllOrders godoc
// @Summary Получить список сделок с фильтрацией
// @Description Возвращает список сделок с возможностью фильтрации, сортировки и пагинации
// @Tags orders
// @Security BearerAuth 
// @Accept json
// @Produce json
// @Param trader_id             query string  false "ID трейдера"
// @Param merchant_id           query string  false "ID мерчанта"
// @Param order_id              query string  false "ID сделки"
// @Param merchant_order_id     query string  false "ID заказа мерчанта"
// @Param status                query string  false "Статус сделки"
// @Param bank_code             query string  false "Код банка"
// @Param time_opening_start    query string  false "Начало периода создания (RFC3339)"
// @Param time_opening_end      query string  false "Конец периода создания (RFC3339)"
// @Param amount_min            query number  false "Минимальная сумма"
// @Param amount_max            query number  false "Максимальная сумма"
// @Param type                  query string  false "Тип сделки"
// @Param payment_system        query string  false "Платежная система"
// @Param device_id             query string  false "ID устройства"
// @Param page                  query int     true  "Номер страницы" default(1)
// @Param limit                 query int     true  "Лимит на страницу" default(50)
// @Param sort                  query string  false "Поле сортировки (amount_fiat, created_at, expires_at) и направление (ASC/DESC), например: amount_fiat DESC"
// @Success 200 {object} orderResponse.GetAllOrdersResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/all [get]
func (h *OrderHandler) GetAllOrders(c *gin.Context) {
    // Парсим параметры запроса
    request, err := parseGetAllOrdersRequest(c)
    if err != nil {
        common.RespondWithError(c, http.StatusBadRequest, err.Error())
        return
    }

    // Вызываем gRPC сервис
    grpcResponse, err := h.OrderClient.GetAllOrders(request)
    if err != nil {
        statusCode := common.GrpcCodeToHTTP(status.Code(err))
        common.RespondWithError(c, statusCode, fmt.Sprintf("Order service error: %v", err))
        return
    }

	response := transformResponse(grpcResponse)

    c.JSON(http.StatusOK, response)
}

// Преобразование gRPC ответа в нашу структуру
func transformResponse(grpcRes *orderpb.GetAllOrdersResponse) *orderResponse.GetAllOrdersResponse {
    response := &orderResponse.GetAllOrdersResponse{
        Orders:     make([]orderResponse.Order, len(grpcRes.Orders)),
        Pagination: transformPagination(grpcRes.Pagination),
    }
    
    for i, grpcOrder := range grpcRes.Orders {
        response.Orders[i] = transformOrder(grpcOrder)
    }
    
    return response
}

func transformPagination(grpcPagination *orderpb.Pagination) orderResponse.Pagination {
    return orderResponse.Pagination{
        CurrentPage:  grpcPagination.CurrentPage,
        TotalPages:   grpcPagination.TotalPages,
        TotalItems:   grpcPagination.TotalItems,
        ItemsPerPage: grpcPagination.ItemsPerPage,
    }
}

func transformOrder(grpcOrder *orderpb.Order) orderResponse.Order {
    // Преобразуем банковские детали
    bankDetail := orderResponse.BankDetail{}
    if grpcOrder.BankDetail != nil {
        bd := grpcOrder.BankDetail
        bankDetail = orderResponse.BankDetail{
            ID:            bd.BankDetailId,
            TraderID:      bd.TraderId,
            Currency:      bd.Currency,
            Country:       bd.Country,
            MinAmount:     bd.MinAmount,
            MaxAmount:     bd.MaxAmount,
            BankName:      bd.BankName,
            PaymentSystem: bd.PaymentSystem,
            Enabled:       bd.Enabled,
            Owner:         bd.Owner,
            CardNumber:    bd.CardNumber,
            Phone:         bd.Phone,
        }
        
        // Преобразуем задержку в строку
        if bd.Delay != nil {
            bankDetail.Delay = bd.Delay.AsDuration().String()
        }
    }
    
    // Преобразуем временные метки
    var expiresAt, createdAt, updatedAt time.Time
    if grpcOrder.ExpiresAt != nil {
        expiresAt = grpcOrder.ExpiresAt.AsTime()
    }
    if grpcOrder.CreatedAt != nil {
        createdAt = grpcOrder.CreatedAt.AsTime()
    }
    if grpcOrder.UpdatedAt != nil {
        updatedAt = grpcOrder.UpdatedAt.AsTime()
    }
    
    return orderResponse.Order{
        OrderID:         grpcOrder.OrderId,
        Status:          grpcOrder.Status,
        AmountFiat:      grpcOrder.AmountFiat,
        AmountCrypto:    grpcOrder.AmountCrypto,
        ExpiresAt:       expiresAt,
        BankDetail:      bankDetail,
        TraderReward:    grpcOrder.TraderRewardPercent,
        CryptoRubRate:   grpcOrder.CryptoRubRate,
        CreatedAt:       createdAt,
        UpdatedAt:       updatedAt,
        MerchantID:      grpcOrder.MerchantId,
        MerchantOrderID: grpcOrder.MerchantOrderId,
    }
}

// Вспомогательная функция для создания строкового указателя
func strPtr(s string) *string {
    if s == "" {
        return nil
    }
    return &s
}

// Вспомогательная функция для парсинга параметров
func parseGetAllOrdersRequest(c *gin.Context) (*orderpb.GetAllOrdersRequest, error) {
    // Парсим основные параметры
    page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
    if err != nil || page < 1 {
        return nil, fmt.Errorf("invalid page parameter")
    }

    limit, err := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 32)
    if err != nil || limit < 1 || limit > 100 {
        return nil, fmt.Errorf("limit must be between 1 and 100")
    }

    // Для float полей используем указатели
    var amountMin, amountMax *float64
    if min := c.Query("amount_min"); min != "" {
        if val, err := strconv.ParseFloat(min, 64); err == nil {
            amountMin = &val
        }
    }
    if max := c.Query("amount_max"); max != "" {
        if val, err := strconv.ParseFloat(max, 64); err == nil {
            amountMax = &val
        }
    }

    // Парсим временные параметры
    var timeStart, timeEnd *timestamppb.Timestamp
    if ts := c.Query("time_opening_start"); ts != "" {
        t, err := time.Parse(time.RFC3339, ts)
        if err != nil {
            return nil, fmt.Errorf("invalid time_opening_start format")
        }
        timeStart = timestamppb.New(t)
    }

    if ts := c.Query("time_opening_end"); ts != "" {
        t, err := time.Parse(time.RFC3339, ts)
        if err != nil {
            return nil, fmt.Errorf("invalid time_opening_end format")
        }
        timeEnd = timestamppb.New(t)
    }

    // Создаем запрос
    req := &orderpb.GetAllOrdersRequest{
        Page:              int32(page),
        Limit:             int32(limit),
        AmountMin:         amountMin,  // теперь *float64
        AmountMax:         amountMax,  // теперь *float64
        TimeOpeningStart:  timeStart,
        TimeOpeningEnd:    timeEnd,
    }

    // Устанавливаем строковые параметры как указатели
    req.TraderId = strPtr(c.Query("trader_id"))
    req.MerchantId = strPtr(c.Query("merchant_id"))
    req.OrderId = strPtr(c.Query("order_id"))
    req.MerchantOrderId = strPtr(c.Query("merchant_order_id"))
    req.Status = strPtr(c.Query("status"))
    req.BankCode = strPtr(c.Query("bank_code"))
    req.Type = strPtr(c.Query("type"))
    req.PaymentSystem = strPtr(c.Query("payment_system"))
    req.DeviceId = strPtr(c.Query("device_id"))
    req.Sort = strPtr(c.Query("sort"))

    return req, nil
}