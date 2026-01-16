package handlers

import (
	"net/http"
	"time"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	orderservice "github.com/LavaJover/shvark-api-gateway/internal/client/order-service"
	orderpb "github.com/LavaJover/shvark-order-service/proto/gen/order"
)

type StoreHandler struct {
	orderClient *orderservice.OrderClient
}

func NewStoreHandler(
	orderClient *orderservice.OrderClient,
) *StoreHandler {
	return &StoreHandler{
		orderClient: orderClient,
	}
}

// Helper functions for conversion
func parseDuration(durationStr string) (*durationpb.Duration, error) {
	if durationStr == "" {
		return nil, nil
	}
	d, err := time.ParseDuration(durationStr)
	if err != nil {
		return nil, err
	}
	return durationpb.New(d), nil
}

func stringToTimePtr(str string) (*time.Time, error) {
	if str == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func storeToHTTP(store *orderpb.Store) StoreHTTP {
	var deletedAt *time.Time
	if store.DeletedAt != nil {
		t := store.DeletedAt.AsTime()
		deletedAt = &t
	}
	
	var merchantInfo MerchantUserInfoHTTP
	if store.MerchantInfo != nil {
		merchantInfo = MerchantUserInfoHTTP{
			MerchantID:       store.MerchantInfo.MerchantId,
			MerchantUsername: store.MerchantInfo.MerchantUsername,
		}
	}
	
	var businessParams StoreBusinessParamsHTTP
	if store.BusinessParams != nil {
		businessParams = StoreBusinessParamsHTTP{
			TrafficType:         FromProtoTrafficType(store.BusinessParams.TrafficType),
			DealPendingDuration: store.BusinessParams.DealPendingDuration.AsDuration().String(),
			DealCreatedDuration: store.BusinessParams.DealCreatedDuration.AsDuration().String(),
			PlatformFee:         store.BusinessParams.PlatformFee,
			Name:                store.BusinessParams.Name,
		}
	}
	
	return StoreHTTP{
		ID:           store.Id,
		MerchantID:   store.MerchantId,
		Enabled:      store.Enabled,
		MerchantInfo: merchantInfo,
		BusinessParams: businessParams,
		CreatedAt:    store.CreatedAt.AsTime(),
		UpdatedAt:    store.UpdatedAt.AsTime(),
		DeletedAt:    deletedAt,
	}
}

func storesToHTTP(stores []*orderpb.Store) []StoreHTTP {
	result := make([]StoreHTTP, len(stores))
	for i, store := range stores {
		result[i] = storeToHTTP(store)
	}
	return result
}

func metricsToHTTP(metrics *orderpb.StoreMetrics) StoreMetricsHTTP {
	return StoreMetricsHTTP{
		StoreID:               metrics.StoreId,
		TotalTraffics:         metrics.TotalTraffics,
		ActiveTraffics:        metrics.ActiveTraffics,
		AverageTraderReward:   metrics.AverageTraderReward,
		UnlockedTrafficsCount: metrics.UnlockedTrafficsCount,
		LockedTrafficsCount:   metrics.LockedTrafficsCount,
		CreatedAt:             metrics.CreatedAt.AsTime(),
		LastActivityAt:        metrics.LastActivityAt.AsTime(),
		TotalRevenue:          metrics.TotalRevenue,
		AverageDealAmount:     metrics.AverageDealAmount,
		TotalDeals:            metrics.TotalDeals,
		SuccessfulDeals:       metrics.SuccessfulDeals,
		FailedDeals:           metrics.FailedDeals,
	}
}

// CreateStore godoc
// @Summary Создание нового магазина
// @Description Создает новый магазин с указанными параметрами
// @Tags store
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateStoreRequestHTTP true "Параметры создания магазина"
// @Success 200 {object} CreateStoreResponseHTTP
// @Failure 400 {object} ErrorResponseHTTP
// @Failure 500 {object} ErrorResponseHTTP
// @Router /stores [post]
func (h *StoreHandler) CreateStore(c *gin.Context) {
	var reqHTTP CreateStoreRequestHTTP
	
	if err := c.ShouldBindJSON(&reqHTTP); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Неверный формат запроса",
			Details: err.Error(),
		})
		return
	}

	dealPendingDuration, err := parseDuration(reqHTTP.DealPendingDuration)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Неверный формат deal_pending_duration",
			Details: err.Error(),
		})
		return
	}

	dealCreatedDuration, err := parseDuration(reqHTTP.DealCreatedDuration)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Неверный формат deal_created_duration",
			Details: err.Error(),
		})
		return
	}

	req := &orderpb.CreateStoreRequest{
		MerchantId:          reqHTTP.MerchantID,
		Enabled:             reqHTTP.Enabled,
		MerchantUsername:    reqHTTP.MerchantUsername,
		TrafficType:         reqHTTP.TrafficType.ToProto(),
		DealPendingDuration: dealPendingDuration,
		DealCreatedDuration: dealCreatedDuration,
		PlatformFee:         reqHTTP.PlatformFee,
		Name:                reqHTTP.Name,
	}

	resp, err := h.orderClient.CreateStore(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponseHTTP{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при создании магазина",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CreateStoreResponseHTTP{
		Store: storeToHTTP(resp.Store),
	})
}

// GetStore godoc
// @Summary Получение информации о магазине
// @Description Возвращает полную информацию о магазине по ID
// @Tags store
// @Produce json
// @Security BearerAuth
// @Param store_id path string true "ID магазина"
// @Success 200 {object} GetStoreResponseHTTP
// @Failure 400 {object} ErrorResponseHTTP
// @Failure 404 {object} ErrorResponseHTTP
// @Failure 500 {object} ErrorResponseHTTP
// @Router /stores/{store_id} [get]
func (h *StoreHandler) GetStore(c *gin.Context) {
	storeID := c.Param("store_id")
	if storeID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Не указан ID магазина",
		})
		return
	}

	req := &orderpb.GetStoreRequest{
		StoreId: storeID,
	}

	resp, err := h.orderClient.GetStore(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponseHTTP{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при получении магазина",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GetStoreResponseHTTP{
		Store: storeToHTTP(resp.Store),
	})
}

// UpdateStore godoc
// @Summary Обновление информации о магазине
// @Description Обновляет параметры существующего магазина
// @Tags store
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param store_id path string true "ID магазина"
// @Param request body UpdateStoreRequestHTTP true "Параметры обновления"
// @Success 200 {object} UpdateStoreResponseHTTP
// @Failure 400 {object} ErrorResponseHTTP
// @Failure 404 {object} ErrorResponseHTTP
// @Failure 500 {object} ErrorResponseHTTP
// @Router /stores/{store_id} [put]
func (h *StoreHandler) UpdateStore(c *gin.Context) {
	storeID := c.Param("store_id")
	if storeID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Не указан ID магазина",
		})
		return
	}

	var reqHTTP UpdateStoreRequestHTTP
	
	if err := c.ShouldBindJSON(&reqHTTP); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Неверный формат запроса",
			Details: err.Error(),
		})
		return
	}

	req := &orderpb.UpdateStoreRequest{
		StoreId: storeID,
		Name:    reqHTTP.Name,
	}

	// Обработка optional полей
	if reqHTTP.Enabled != nil {
		req.Enabled = wrapperspb.Bool(*reqHTTP.Enabled)
	}
	
	if reqHTTP.MerchantUsername != "" {
		req.MerchantUsername = reqHTTP.MerchantUsername
	}
	
	if reqHTTP.TrafficType != nil {
		trafficType := (*reqHTTP.TrafficType).ToProto()
		req.TrafficType = trafficType
	}
	
	if reqHTTP.DealPendingDuration != nil {
		duration, err := parseDuration(*reqHTTP.DealPendingDuration)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
				Code:    http.StatusBadRequest,
				Message: "Неверный формат deal_pending_duration",
				Details: err.Error(),
			})
			return
		}
		req.DealPendingDuration = duration
	}
	
	if reqHTTP.DealCreatedDuration != nil {
		duration, err := parseDuration(*reqHTTP.DealCreatedDuration)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
				Code:    http.StatusBadRequest,
				Message: "Неверный формат deal_created_duration",
				Details: err.Error(),
			})
			return
		}
		req.DealCreatedDuration = duration
	}
	
	if reqHTTP.PlatformFee != nil {
		req.PlatformFee = wrapperspb.Double(*reqHTTP.PlatformFee)
	}

	resp, err := h.orderClient.UpdateStore(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponseHTTP{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при обновлении магазина",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, UpdateStoreResponseHTTP{
		Store: storeToHTTP(resp.Store),
	})
}

// DeleteStore godoc
// @Summary Удаление магазина
// @Description Удаляет магазин по ID с возможностью принудительного удаления
// @Tags store
// @Produce json
// @Security BearerAuth
// @Param store_id path string true "ID магазина"
// @Param force query boolean false "Принудительное удаление" default(false)
// @Success 200 {object} DeleteStoreResponseHTTP
// @Failure 400 {object} ErrorResponseHTTP
// @Failure 404 {object} ErrorResponseHTTP
// @Failure 500 {object} ErrorResponseHTTP
// @Router /stores/{store_id} [delete]
func (h *StoreHandler) DeleteStore(c *gin.Context) {
	storeID := c.Param("store_id")
	if storeID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Не указан ID магазина",
		})
		return
	}

	forceStr := c.DefaultQuery("force", "false")
	force, err := strconv.ParseBool(forceStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Неверный формат параметра force",
			Details: err.Error(),
		})
		return
	}

	req := &orderpb.DeleteStoreRequest{
		StoreId: storeID,
		Force:   force,
	}

	resp, err := h.orderClient.DeleteStore(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponseHTTP{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при удалении магазина",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, DeleteStoreResponseHTTP{
		Success: resp.Success,
		Message: resp.Message,
	})
}

// ListStores godoc
// @Summary Список магазинов
// @Description Возвращает список магазинов с пагинацией и фильтрацией
// @Tags store
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ListStoresRequestHTTP true "Параметры пагинации и фильтрации"
// @Success 200 {object} ListStoresResponseHTTP
// @Failure 400 {object} ErrorResponseHTTP
// @Failure 500 {object} ErrorResponseHTTP
// @Router /stores/list [post]
func (h *StoreHandler) ListStores(c *gin.Context) {
	var reqHTTP ListStoresRequestHTTP
	
	if err := c.ShouldBindJSON(&reqHTTP); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Неверный формат запроса",
			Details: err.Error(),
		})
		return
	}

	// Преобразование HTTP -> gRPC
	req := &orderpb.ListStoresRequest{
		Pagination: &orderpb.PaginationParams{
			Page:  reqHTTP.Pagination.Page,
			Limit: reqHTTP.Pagination.Limit,
		},
		Filters: &orderpb.StoreFilters{
			MerchantId:  reqHTTP.Filters.MerchantID,
			Name:        reqHTTP.Filters.Name,
			TrafficType: reqHTTP.Filters.TrafficType.ToProto(),
			StoreIds:    reqHTTP.Filters.StoreIDs,
		},
	}

	// Обработка optional полей
	if reqHTTP.Filters.Enabled != nil {
		req.Filters.Enabled = wrapperspb.Bool(*reqHTTP.Filters.Enabled)
	}
	
	if reqHTTP.Filters.CreatedFrom != nil {
		req.Filters.CreatedFrom = timestamppb.New(*reqHTTP.Filters.CreatedFrom)
	}
	
	if reqHTTP.Filters.CreatedTo != nil {
		req.Filters.CreatedTo = timestamppb.New(*reqHTTP.Filters.CreatedTo)
	}

	resp, err := h.orderClient.ListStores(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponseHTTP{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при получении списка магазинов",
			Details: err.Error(),
		})
		return
	}

	// Преобразование gRPC -> HTTP
	c.JSON(http.StatusOK, ListStoresResponseHTTP{
		Stores: storesToHTTP(resp.Stores),
		PaginationInfo: PaginationInfoHTTP{
			CurrentPage: resp.PaginationInfo.CurrentPage,
			Limit:       resp.PaginationInfo.Limit,
			TotalItems:  resp.PaginationInfo.TotalItems,
			TotalPages:  resp.PaginationInfo.TotalPages,
			HasNext:     resp.PaginationInfo.HasNext,
			HasPrev:     resp.PaginationInfo.HasPrev,
		},
	})
}

// GetStoreWithTraffics godoc
// @Summary Получение магазина с трафиками
// @Description Возвращает информацию о магазине вместе с его трафиками
// @Tags store
// @Produce json
// @Security BearerAuth
// @Param store_id path string true "ID магазина"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Лимит записей" default(20)
// @Param trader_filter query string false "Фильтр по трейдеру"
// @Param include_traffic_details query boolean false "Включить детали трафика" default(false)
// @Success 200 {object} GetStoreWithTrafficsResponseHTTP
// @Failure 400 {object} ErrorResponseHTTP
// @Failure 404 {object} ErrorResponseHTTP
// @Failure 500 {object} ErrorResponseHTTP
// @Router /stores/{store_id}/traffics [get]
func (h *StoreHandler) GetStoreWithTraffics(c *gin.Context) {
	storeID := c.Param("store_id")
	if storeID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Не указан ID магазина",
		})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")
	traderFilter := c.Query("trader_filter")
	includeTrafficDetailsStr := c.DefaultQuery("include_traffic_details", "false")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Неверный формат параметра page",
			Details: err.Error(),
		})
		return
	}

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Неверный формат параметра limit",
			Details: err.Error(),
		})
		return
	}

	includeTrafficDetails, err := strconv.ParseBool(includeTrafficDetailsStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Неверный формат параметра include_traffic_details",
			Details: err.Error(),
		})
		return
	}

	req := &orderpb.GetStoreWithTrafficsRequest{
		StoreId:                storeID,
		Pagination:             &orderpb.PaginationParams{
			Page:  int32(page),
			Limit: int32(limit),
		},
		TraderFilter:           traderFilter,
		IncludeTrafficDetails:  includeTrafficDetails,
	}

	resp, err := h.orderClient.GetStoreWithTraffics(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponseHTTP{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при получении магазина с трафиками",
			Details: err.Error(),
		})
		return
	}

	// Преобразование TrafficRecords (упрощенное - пока не реализовано полностью)
	trafficRecords := make([]TrafficRecordHTTP, len(resp.TrafficRecords))
	for i, tr := range resp.TrafficRecords {
		trafficRecords[i] = TrafficRecordHTTP{
			ID:                  tr.Id,
			StoreID:             tr.StoreId,
			TraderID:            tr.TraderId,
			TraderRewardPercent: tr.TraderRewardPercent,
			TraderPriority:      tr.TraderPriority,
			CreatedAt:           tr.CreatedAt.AsTime(),
			UpdatedAt:           tr.UpdatedAt.AsTime(),
		}
	}

	c.JSON(http.StatusOK, GetStoreWithTrafficsResponseHTTP{
		Store:         storeToHTTP(resp.Store),
		TrafficRecords: trafficRecords,
		PaginationInfo: PaginationInfoHTTP{
			CurrentPage: resp.PaginationInfo.CurrentPage,
			Limit:       resp.PaginationInfo.Limit,
			TotalItems:  resp.PaginationInfo.TotalItems,
			TotalPages:  resp.PaginationInfo.TotalPages,
			HasNext:     resp.PaginationInfo.HasNext,
			HasPrev:     resp.PaginationInfo.HasPrev,
		},
		TotalTraffics: resp.TotalTraffics,
		TraderFilter:  resp.TraderFilter,
	})
}

// GetStoreByTrafficId godoc
// @Summary Получение магазина по ID трафика
// @Description Возвращает информацию о магазине, к которому относится указанный трафик
// @Tags store
// @Produce json
// @Security BearerAuth
// @Param traffic_id path string true "ID трафика"
// @Success 200 {object} GetStoreByTrafficIdResponseHTTP
// @Failure 400 {object} ErrorResponseHTTP
// @Failure 404 {object} ErrorResponseHTTP
// @Failure 500 {object} ErrorResponseHTTP
// @Router /traffics/{traffic_id}/store [get]
func (h *StoreHandler) GetStoreByTrafficId(c *gin.Context) {
	trafficID := c.Param("traffic_id")
	if trafficID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Не указан ID трафика",
		})
		return
	}

	req := &orderpb.GetStoreByTrafficIdRequest{
		TrafficId: trafficID,
	}

	resp, err := h.orderClient.GetStoreByTrafficId(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponseHTTP{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при получении магазина по трафику",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GetStoreByTrafficIdResponseHTTP{
		Store: storeToHTTP(resp.Store),
	})
}

// CheckStoreNameUnique godoc
// @Summary Проверка уникальности имени магазина
// @Description Проверяет, уникально ли имя магазина для указанного мерчанта
// @Tags store
// @Produce json
// @Security BearerAuth
// @Param merchant_id query string true "ID мерчанта"
// @Param store_name query string true "Имя магазина"
// @Success 200 {object} CheckStoreNameUniqueResponseHTTP
// @Failure 400 {object} ErrorResponseHTTP
// @Failure 500 {object} ErrorResponseHTTP
// @Router /stores/check-name [get]
func (h *StoreHandler) CheckStoreNameUnique(c *gin.Context) {
	merchantID := c.Query("merchant_id")
	storeName := c.Query("store_name")
	
	if merchantID == "" || storeName == "" {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Не указаны merchant_id или store_name",
		})
		return
	}

	req := &orderpb.CheckStoreNameUniqueRequest{
		MerchantId: merchantID,
		StoreName:  storeName,
	}

	resp, err := h.orderClient.CheckStoreNameUnique(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponseHTTP{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при проверке уникальности имени магазина",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CheckStoreNameUniqueResponseHTTP{
		IsUnique:   resp.IsUnique,
		Message:    resp.Message,
		MerchantID: resp.MerchantId,
		StoreName:  resp.StoreName,
	})
}

// ValidateStoreForTraffic godoc
// @Summary Валидация магазина для трафика
// @Description Проверяет, готов ли магазин к приему трафика
// @Tags store
// @Produce json
// @Security BearerAuth
// @Param store_id path string true "ID магазина"
// @Success 200 {object} ValidateStoreForTrafficResponseHTTP
// @Failure 400 {object} ErrorResponseHTTP
// @Failure 404 {object} ErrorResponseHTTP
// @Failure 500 {object} ErrorResponseHTTP
// @Router /stores/{store_id}/validate-traffic [get]
func (h *StoreHandler) ValidateStoreForTraffic(c *gin.Context) {
	storeID := c.Param("store_id")
	if storeID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Не указан ID магазина",
		})
		return
	}

	req := &orderpb.ValidateStoreForTrafficRequest{
		StoreId: storeID,
	}

	resp, err := h.orderClient.ValidateStoreForTraffic(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponseHTTP{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при валидации магазина для трафика",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ValidateStoreForTrafficResponseHTTP{
		Valid:   resp.Valid,
		Message: resp.Message,
		Errors:  resp.Errors,
	})
}

// ToggleStoreStatus godoc
// @Summary Переключение статуса магазина
// @Description Включает или выключает магазин
// @Tags store
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param store_id path string true "ID магазина"
// @Param request body ToggleStoreStatusRequestHTTP true "Параметры переключения статуса"
// @Success 200 {object} ToggleStoreStatusResponseHTTP
// @Failure 400 {object} ErrorResponseHTTP
// @Failure 404 {object} ErrorResponseHTTP
// @Failure 500 {object} ErrorResponseHTTP
// @Router /stores/{store_id}/toggle-status [post]
func (h *StoreHandler) ToggleStoreStatus(c *gin.Context) {
	storeID := c.Param("store_id")
	if storeID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Не указан ID магазина",
		})
		return
	}

	var reqHTTP ToggleStoreStatusRequestHTTP
	
	if err := c.ShouldBindJSON(&reqHTTP); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Неверный формат запроса",
			Details: err.Error(),
		})
		return
	}

	req := &orderpb.ToggleStoreStatusRequest{
		StoreId: storeID,
		Enabled: reqHTTP.Enabled,
	}

	resp, err := h.orderClient.ToggleStoreStatus(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponseHTTP{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при переключении статуса магазина",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ToggleStoreStatusResponseHTTP{
		Store:         storeToHTTP(resp.Store),
		PreviousStatus: resp.PreviousStatus,
		NewStatus:      resp.NewStatus,
	})
}

// EnableStore godoc
// @Summary Включение магазина
// @Description Включает магазин (устанавливает enabled=true)
// @Tags store
// @Produce json
// @Security BearerAuth
// @Param store_id path string true "ID магазина"
// @Success 200 {object} EnableStoreResponseHTTP
// @Failure 400 {object} ErrorResponseHTTP
// @Failure 404 {object} ErrorResponseHTTP
// @Failure 500 {object} ErrorResponseHTTP
// @Router /stores/{store_id}/enable [post]
func (h *StoreHandler) EnableStore(c *gin.Context) {
	storeID := c.Param("store_id")
	if storeID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Не указан ID магазина",
		})
		return
	}

	req := &orderpb.EnableStoreRequest{
		StoreId: storeID,
	}

	resp, err := h.orderClient.EnableStore(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponseHTTP{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при включении магазина",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, EnableStoreResponseHTTP{
		Store: storeToHTTP(resp.Store),
	})
}

// DisableStore godoc
// @Summary Отключение магазина
// @Description Отключает магазин (устанавливает enabled=false)
// @Tags store
// @Produce json
// @Security BearerAuth
// @Param store_id path string true "ID магазина"
// @Param force query boolean false "Принудительное отключение" default(false)
// @Success 200 {object} DisableStoreResponseHTTP
// @Failure 400 {object} ErrorResponseHTTP
// @Failure 404 {object} ErrorResponseHTTP
// @Failure 500 {object} ErrorResponseHTTP
// @Router /stores/{store_id}/disable [post]
func (h *StoreHandler) DisableStore(c *gin.Context) {
	storeID := c.Param("store_id")
	if storeID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Не указан ID магазина",
		})
		return
	}

	forceStr := c.DefaultQuery("force", "false")
	force, err := strconv.ParseBool(forceStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Неверный формат параметра force",
			Details: err.Error(),
		})
		return
	}

	req := &orderpb.DisableStoreRequest{
		StoreId: storeID,
		Force:   force,
	}

	resp, err := h.orderClient.DisableStore(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponseHTTP{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при отключении магазина",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, DisableStoreResponseHTTP{
		Store: storeToHTTP(resp.Store),
	})
}

// BulkUpdateStoresStatus godoc
// @Summary Массовое обновление статуса магазинов
// @Description Включает или выключает несколько магазинов одновременно
// @Tags store
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body BulkUpdateStoresStatusRequestHTTP true "Параметры массового обновления"
// @Success 200 {object} BulkUpdateStoresStatusResponseHTTP
// @Failure 400 {object} ErrorResponseHTTP
// @Failure 500 {object} ErrorResponseHTTP
// @Router /stores/bulk-update-status [post]
func (h *StoreHandler) BulkUpdateStoresStatus(c *gin.Context) {
	var reqHTTP BulkUpdateStoresStatusRequestHTTP
	
	if err := c.ShouldBindJSON(&reqHTTP); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Неверный формат запроса",
			Details: err.Error(),
		})
		return
	}

	req := &orderpb.BulkUpdateStoresStatusRequest{
		StoreIds: reqHTTP.StoreIDs,
		Enabled:  reqHTTP.Enabled,
	}

	resp, err := h.orderClient.BulkUpdateStoresStatus(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponseHTTP{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при массовом обновлении статуса магазинов",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, BulkUpdateStoresStatusResponseHTTP{
		UpdatedCount:   resp.UpdatedCount,
		FailedCount:    resp.FailedCount,
		FailedStoreIDs: resp.FailedStoreIds,
		Errors:         resp.Errors,
	})
}

// SearchStores godoc
// @Summary Поиск магазинов
// @Description Поиск магазинов с расширенными фильтрами, сортировкой и пагинацией
// @Tags store
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SearchStoresRequestHTTP true "Параметры поиска"
// @Success 200 {object} SearchStoresResponseHTTP
// @Failure 400 {object} ErrorResponseHTTP
// @Failure 500 {object} ErrorResponseHTTP
// @Router /stores/search [post]
func (h *StoreHandler) SearchStores(c *gin.Context) {
	var reqHTTP SearchStoresRequestHTTP
	
	if err := c.ShouldBindJSON(&reqHTTP); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Неверный формат запроса",
			Details: err.Error(),
		})
		return
	}

	// Преобразование HTTP -> gRPC
	req := &orderpb.SearchStoresRequest{
		Pagination: &orderpb.PaginationParams{
			Page:  reqHTTP.Pagination.Page,
			Limit: reqHTTP.Pagination.Limit,
		},
		Filters: &orderpb.StoreFilters{
			MerchantId:  reqHTTP.Filters.MerchantID,
			Name:        reqHTTP.Filters.Name,
			TrafficType: reqHTTP.Filters.TrafficType.ToProto(),
			StoreIds:    reqHTTP.Filters.StoreIDs,
		},
		SearchQuery: reqHTTP.SearchQuery,
		SortBy:      reqHTTP.SortBy,
		SortDesc:    reqHTTP.SortDesc,
	}

	// Обработка optional полей
	if reqHTTP.Filters.Enabled != nil {
		req.Filters.Enabled = wrapperspb.Bool(*reqHTTP.Filters.Enabled)
	}
	
	if reqHTTP.Filters.CreatedFrom != nil {
		req.Filters.CreatedFrom = timestamppb.New(*reqHTTP.Filters.CreatedFrom)
	}
	
	if reqHTTP.Filters.CreatedTo != nil {
		req.Filters.CreatedTo = timestamppb.New(*reqHTTP.Filters.CreatedTo)
	}

	resp, err := h.orderClient.SearchStores(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponseHTTP{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при поиске магазинов",
			Details: err.Error(),
		})
		return
	}

	// Преобразование gRPC -> HTTP
	c.JSON(http.StatusOK, SearchStoresResponseHTTP{
		Stores: storesToHTTP(resp.Stores),
		PaginationInfo: PaginationInfoHTTP{
			CurrentPage: resp.PaginationInfo.CurrentPage,
			Limit:       resp.PaginationInfo.Limit,
			TotalItems:  resp.PaginationInfo.TotalItems,
			TotalPages:  resp.PaginationInfo.TotalPages,
			HasNext:     resp.PaginationInfo.HasNext,
			HasPrev:     resp.PaginationInfo.HasPrev,
		},
		AppliedFilters: StoreFiltersHTTP{
			MerchantID:  resp.AppliedFilters.MerchantId,
			Name:        resp.AppliedFilters.Name,
			TrafficType: FromProtoTrafficType(resp.AppliedFilters.TrafficType),
			StoreIDs:    resp.AppliedFilters.StoreIds,
		},
		SearchQuery: resp.SearchQuery,
	})
}

// GetStoresByMerchant godoc
// @Summary Получение магазинов мерчанта
// @Description Возвращает список магазинов, принадлежащих указанному мерчанту
// @Tags store
// @Produce json
// @Security BearerAuth
// @Param merchant_id path string true "ID мерчанта"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Лимит записей" default(20)
// @Param only_active query boolean false "Только активные магазины" default(false)
// @Success 200 {object} GetStoresByMerchantResponseHTTP
// @Failure 400 {object} ErrorResponseHTTP
// @Failure 404 {object} ErrorResponseHTTP
// @Failure 500 {object} ErrorResponseHTTP
// @Router /merchants/{merchant_id}/stores [get]
func (h *StoreHandler) GetStoresByMerchant(c *gin.Context) {
	merchantID := c.Param("merchant_id")
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Не указан ID мерчанта",
		})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")
	onlyActiveStr := c.DefaultQuery("only_active", "false")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Неверный формат параметра page",
			Details: err.Error(),
		})
		return
	}

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Неверный формат параметра limit",
			Details: err.Error(),
		})
		return
	}

	onlyActive, err := strconv.ParseBool(onlyActiveStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Неверный формат параметра only_active",
			Details: err.Error(),
		})
		return
	}

	req := &orderpb.GetStoresByMerchantRequest{
		MerchantId: merchantID,
		Pagination: &orderpb.PaginationParams{
			Page:  int32(page),
			Limit: int32(limit),
		},
		OnlyActive: onlyActive,
	}

	resp, err := h.orderClient.GetStoresByMerchant(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponseHTTP{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при получении магазинов мерчанта",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GetStoresByMerchantResponseHTTP{
		Stores: storesToHTTP(resp.Stores),
		PaginationInfo: PaginationInfoHTTP{
			CurrentPage: resp.PaginationInfo.CurrentPage,
			Limit:       resp.PaginationInfo.Limit,
			TotalItems:  resp.PaginationInfo.TotalItems,
			TotalPages:  resp.PaginationInfo.TotalPages,
			HasNext:     resp.PaginationInfo.HasNext,
			HasPrev:     resp.PaginationInfo.HasPrev,
		},
		MerchantID:   resp.MerchantId,
		TotalStores:  resp.TotalStores,
		ActiveStores: resp.ActiveStores,
	})
}

// GetActiveStores godoc
// @Summary Получение активных магазинов
// @Description Возвращает список активных магазинов мерчанта
// @Tags store
// @Produce json
// @Security BearerAuth
// @Param merchant_id query string true "ID мерчанта"
// @Success 200 {object} GetActiveStoresResponseHTTP
// @Failure 400 {object} ErrorResponseHTTP
// @Failure 500 {object} ErrorResponseHTTP
// @Router /stores/active [get]
func (h *StoreHandler) GetActiveStores(c *gin.Context) {
	merchantID := c.Query("merchant_id")
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Не указан ID мерчанта",
		})
		return
	}

	req := &orderpb.GetActiveStoresRequest{
		MerchantId: merchantID,
	}

	resp, err := h.orderClient.GetActiveStores(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponseHTTP{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при получении активных магазинов",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GetActiveStoresResponseHTTP{
		Stores:     storesToHTTP(resp.Stores),
		MerchantID: resp.MerchantId,
		Count:      resp.Count,
	})
}

// GetStoreMetrics godoc
// @Summary Получение метрик магазина
// @Description Возвращает метрики магазина за указанный период
// @Tags store
// @Produce json
// @Security BearerAuth
// @Param store_id path string true "ID магазина"
// @Param period_from query string false "Начало периода (RFC3339)"
// @Param period_to query string false "Конец периода (RFC3339)"
// @Success 200 {object} GetStoreMetricsResponseHTTP
// @Failure 400 {object} ErrorResponseHTTP
// @Failure 404 {object} ErrorResponseHTTP
// @Failure 500 {object} ErrorResponseHTTP
// @Router /stores/{store_id}/metrics [get]
func (h *StoreHandler) GetStoreMetrics(c *gin.Context) {
	storeID := c.Param("store_id")
	if storeID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Не указан ID магазина",
		})
		return
	}

	periodFromStr := c.Query("period_from")
	periodToStr := c.Query("period_to")

	req := &orderpb.GetStoreMetricsRequest{
		StoreId: storeID,
	}

	// Парсинг optional параметров периода
	if periodFromStr != "" {
		periodFrom, err := stringToTimePtr(periodFromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
				Code:    http.StatusBadRequest,
				Message: "Неверный формат period_from",
				Details: err.Error(),
			})
			return
		}
		if periodFrom != nil {
			req.PeriodFrom = timestamppb.New(*periodFrom)
		}
	}

	if periodToStr != "" {
		periodTo, err := stringToTimePtr(periodToStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
				Code:    http.StatusBadRequest,
				Message: "Неверный формат period_to",
				Details: err.Error(),
			})
			return
		}
		if periodTo != nil {
			req.PeriodTo = timestamppb.New(*periodTo)
		}
	}

	resp, err := h.orderClient.GetStoreMetrics(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponseHTTP{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при получении метрик магазина",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GetStoreMetricsResponseHTTP{
		Metrics: metricsToHTTP(resp.Metrics),
	})
}

// CalculateStoreMetrics godoc
// @Summary Расчет метрик магазина
// @Description Запускает расчет метрик магазина в реальном времени
// @Tags store
// @Produce json
// @Security BearerAuth
// @Param store_id path string true "ID магазина"
// @Success 200 {object} CalculateStoreMetricsResponseHTTP
// @Failure 400 {object} ErrorResponseHTTP
// @Failure 404 {object} ErrorResponseHTTP
// @Failure 500 {object} ErrorResponseHTTP
// @Router /stores/{store_id}/calculate-metrics [post]
func (h *StoreHandler) CalculateStoreMetrics(c *gin.Context) {
	storeID := c.Param("store_id")
	if storeID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Не указан ID магазина",
		})
		return
	}

	req := &orderpb.CalculateStoreMetricsRequest{
		StoreId: storeID,
	}

	resp, err := h.orderClient.CalculateStoreMetrics(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponseHTTP{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при расчете метрик магазина",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CalculateStoreMetricsResponseHTTP{
		Metrics:      metricsToHTTP(resp.Metrics),
		CalculatedAt: resp.CalculatedAt.AsTime(),
	})
}

// BatchGetStores godoc
// @Summary Пакетное получение магазинов
// @Description Возвращает информацию о нескольких магазинах по их ID
// @Tags store
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body BatchGetStoresRequestHTTP true "Список ID магазинов"
// @Success 200 {object} BatchGetStoresResponseHTTP
// @Failure 400 {object} ErrorResponseHTTP
// @Failure 500 {object} ErrorResponseHTTP
// @Router /stores/batch [post]
func (h *StoreHandler) BatchGetStores(c *gin.Context) {
	var reqHTTP BatchGetStoresRequestHTTP
	
	if err := c.ShouldBindJSON(&reqHTTP); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponseHTTP{
			Code:    http.StatusBadRequest,
			Message: "Неверный формат запроса",
			Details: err.Error(),
		})
		return
	}

	req := &orderpb.BatchGetStoresRequest{
		StoreIds: reqHTTP.StoreIDs,
	}

	resp, err := h.orderClient.BatchGetStores(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponseHTTP{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при пакетном получении магазинов",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, BatchGetStoresResponseHTTP{
		Stores:      storesToHTTP(resp.Stores),
		NotFoundIDs: resp.NotFoundIds,
	})
}

// HealthCheck godoc
// @Summary Проверка здоровья сервиса
// @Description Возвращает статус здоровья сервиса магазинов
// @Tags store
// @Produce json
// @Security BearerAuth
// @Success 200 {object} HealthCheckResponseHTTP
// @Failure 500 {object} ErrorResponseHTTP
// @Router /stores/health [get]
func (h *StoreHandler) HealthCheck(c *gin.Context) {
	resp, err := h.orderClient.HealthCheck()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponseHTTP{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка при проверке здоровья",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, HealthCheckResponseHTTP{
		Status:    resp.Status,
		Version:   resp.Version,
		Timestamp: resp.Timestamp.AsTime(),
		Services:  resp.Services,
	})
}

// Функция для регистрации всех маршрутов
func (h *StoreHandler) RegisterRoutes(router *gin.RouterGroup) {
	storeGroup := router.Group("/stores")
	{
		storeGroup.POST("", h.CreateStore)
		storeGroup.GET("/:store_id", h.GetStore)
		storeGroup.PUT("/:store_id", h.UpdateStore)
		storeGroup.DELETE("/:store_id", h.DeleteStore)
		storeGroup.POST("/list", h.ListStores)
		storeGroup.GET("/:store_id/traffics", h.GetStoreWithTraffics)
		storeGroup.GET("/check-name", h.CheckStoreNameUnique)
		storeGroup.GET("/:store_id/validate-traffic", h.ValidateStoreForTraffic)
		storeGroup.POST("/:store_id/toggle-status", h.ToggleStoreStatus)
		storeGroup.POST("/:store_id/enable", h.EnableStore)
		storeGroup.POST("/:store_id/disable", h.DisableStore)
		storeGroup.POST("/bulk-update-status", h.BulkUpdateStoresStatus)
		storeGroup.POST("/search", h.SearchStores)
		storeGroup.GET("/active", h.GetActiveStores)
		storeGroup.GET("/:store_id/metrics", h.GetStoreMetrics)
		storeGroup.POST("/:store_id/calculate-metrics", h.CalculateStoreMetrics)
		storeGroup.POST("/batch", h.BatchGetStores)
		storeGroup.GET("/health", h.HealthCheck)
	}

	// Дополнительные маршруты для связи с другими сущностями
	trafficGroup := router.Group("/traffics")
	{
		trafficGroup.GET("/:traffic_id/store", h.GetStoreByTrafficId)
	}

	merchantGroup := router.Group("/merchants")
	{
		merchantGroup.GET("/:merchant_id/stores", h.GetStoresByMerchant)
	}
}