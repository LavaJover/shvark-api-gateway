package handlers

import (
    "context"
    "log"
    "net/http"
    "strconv"
    "time"

    "github.com/LavaJover/shvark-api-gateway/internal/client"
    orderpb "github.com/LavaJover/shvark-order-service/proto/gen"
    "github.com/gin-gonic/gin"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type AutomaticHandler struct {
    orderService *client.OrderClient
	deviceService *client.DeviceClient
}

func NewAutomaticHandler(
	orderService *client.OrderClient,
	deviceService *client.DeviceClient,
) *AutomaticHandler {
    return &AutomaticHandler{
        orderService: orderService,
		deviceService: deviceService,
    }
}

// ==================== SMS/PAYMENT PROCESSING ====================

type SMSRequest struct {
    Success       bool     `json:"success"`
    PaymentSystem string   `json:"payment_system"`
    Amount        float64  `json:"amount"`
    Balance       float64  `json:"balance"`
    Group         string   `json:"group"`
    Methods       []string `json:"methods"`
    Direction     string   `json:"direction"`
    Text          string   `json:"text"`
    Title         string   `json:"title"`
    ReceivedAt    int64    `json:"received_at"`
    Version       string   `json:"version"`
    Blocked       bool     `json:"blocked"`
    TooOld        bool     `json:"too_old"`
    Unknown       bool     `json:"unknown"`
}

// Sms обрабатывает входящее СМС уведомление с телефона трейдера
// @Summary Process SMS notification
// @Description Process payment notification from trader's phone
// @Tags automatic
// @Accept json
// @Produce json
// @Param sms body SMSRequest true "SMS data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /automatic/process-sms [post]
func (h *AutomaticHandler) Sms(c *gin.Context) {
    var req SMSRequest
    
    if err := c.BindJSON(&req); err != nil {
        log.Printf("❌ [SMS] Parse error: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

    log.Printf("📱 [SMS] Received SMS: device=%s, amount=%.2f, payment_system=%s, direction=%s",
        req.Group, req.Amount, req.PaymentSystem, req.Direction)

    // Валидация входящего уведомления
    if !h.validateSMS(req) {
        log.Printf("⚠️  [SMS] Validation failed for device=%s: success=%v, blocked=%v, too_old=%v, unknown=%v",
            req.Group, req.Success, req.Blocked, req.TooOld, req.Unknown)
        c.JSON(http.StatusOK, gin.H{
            "status": "ignored",
            "reason": "validation failed",
        })
        return
    }

    // Подготовка данных для gRPC вызова
    grpcReq := &orderpb.ProcessAutomaticPaymentRequest{
        Group:         req.Group,
        Amount:        req.Amount,
        PaymentSystem: req.PaymentSystem,
        Methods:       req.Methods,
        Direction:     req.Direction,
        ReceivedAt:    req.ReceivedAt,
        Text:          req.Text,
        Metadata: map[string]string{
            "title":   req.Title,
            "balance": strconv.FormatFloat(req.Balance, 'f', 2, 64),
        },
    }

    // Вызов order-service с retry логикой
    ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
    defer cancel()

    response, err := h.withRetry(ctx, 3, func() (*orderpb.ProcessAutomaticPaymentResponse, error) {
        return h.orderService.ProcessAutomaticPayment(ctx, grpcReq)
    })

    if err != nil {
        log.Printf("❌ [SMS] Processing error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":  "processing failed",
            "detail": err.Error(),
        })
        return
    }

    log.Printf("✅ [SMS] Processed: device=%s, action=%s, success=%v, orders=%d",
        req.Group, response.Action, response.Success, len(response.Results))

    c.JSON(http.StatusOK, gin.H{
        "status":    "processed",
        "order_id":  response.OrderId,
        "action":    response.Action,
        "processed": response.Success,
        "results":   response.Results,
    })
}

// validateSMS валидирует входящее SMS уведомление
func (h *AutomaticHandler) validateSMS(req SMSRequest) bool {
    // Игнорируем неуспешные уведомления
    if !req.Success || req.Unknown {
        return false
    }

    // Игнорируем заблокированные или устаревшие уведомления
    if req.Blocked || req.TooOld {
        return false
    }

    // Проверяем обязательные поля
    if req.Group == "" || req.Amount <= 0 {
        return false
    }

    // Проверяем направление платежа (только входящие)
    if req.Direction != "in" {
        return false
    }

    return true
}

// withRetry выполняет функцию с повторными попытками
func (h *AutomaticHandler) withRetry(ctx context.Context, maxRetries int, fn func() (*orderpb.ProcessAutomaticPaymentResponse, error)) (*orderpb.ProcessAutomaticPaymentResponse, error) {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        resp, err := fn()
        if err == nil {
            return resp, nil
        }
        
        // Проверяем, стоит ли повторять запрос
        grpcCode := status.Code(err)
        if grpcCode == codes.Unavailable || grpcCode == codes.DeadlineExceeded {
            lastErr = err
            backoff := time.Duration(i*100) * time.Millisecond
            log.Printf("⏳ [SMS] Retry %d/%d after %v: %v", i+1, maxRetries, backoff, err)
            time.Sleep(backoff)
            continue
        }
        
        return nil, err
    }
    
    return nil, lastErr
}

// ==================== DEVICE LIVENESS ====================

// Live обрабатывает пинги от устройств трейдеров (keepalive сигналы)
// @Summary Device liveness ping
// @Description Receive liveness ping from trader's device (heartbeat)
// @Tags automatic
// @Accept json
// @Produce json
// @Param body body map[string]string true "Device group"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /automatic/liveness [post]
func (h *AutomaticHandler) Live(c *gin.Context) {
    var body map[string]interface{}
    
    if err := c.BindJSON(&body); err != nil {
        log.Println("❌ [LIVENESS] Parse error:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }
    
    group, ok := body["group"].(string)
    if !ok || group == "" {
        log.Println("⚠️  [LIVENESS] Missing group field")
        c.JSON(http.StatusBadRequest, gin.H{"error": "group field is required"})
        return
    }
    
    log.Printf("💓 [LIVENESS] Ping received: device=%s", group)
    
    // Вызываем order-service для обновления статуса
    ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
    defer cancel()
    
    _, err := h.orderService.UpdateDeviceLiveness(ctx, &orderpb.UpdateDeviceLivenessRequest{
        DeviceId: group,
    })
    
    if err != nil {
        log.Printf("❌ [LIVENESS] Failed to update device status: %v", err)
        // Не возвращаем ошибку клиенту - он все равно получит подтверждение
    }
    
    c.JSON(http.StatusOK, gin.H{
        "alive":     true,
        "timestamp": time.Now().Unix(),
    })
}

// ==================== DEVICE AUTH ====================

// Auth обрабатывает авторизацию устройства (QR код сканирование)
// @Summary Device authorization
// @Description Authorize device via QR code scan from trader account
// @Tags automatic
// @Accept json
// @Produce json
// @Param body body map[string]string true "Device auth data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /automatic/auth [post]
func (h *AutomaticHandler) Auth(c *gin.Context) {
    var body map[string]interface{}
    
    if err := c.BindJSON(&body); err != nil {
        log.Printf("❌ [AUTH] Parse error: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

    group, ok := body["group"].(string)
    if !ok || group == "" {
        log.Println("⚠️  [AUTH] Missing group field")
        c.JSON(http.StatusBadRequest, gin.H{"error": "group field is required"})
        return
    }

    log.Printf("🔐 [AUTH] Auth request: device=%s", group)

    // Можно добавить проверку токена в заголовке Authorization
    authToken := c.GetHeader("Authorization")
    if authToken == "" {
        log.Printf("⚠️  [AUTH] Missing authorization token")
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":    "authorized",
        "device":    group,
        "timestamp": time.Now().Unix(),
    })
}

// ==================== AUTOMATIC LOGS ====================

// GetAutomaticLogs получает логи автоматической обработки платежей
// @Summary Get automatic payment logs
// @Description Retrieve logs of automatic payment processing with filters
// @Tags automatic
// @Accept json
// @Produce json
// @Param trader_id query string false "Filter by trader ID"
// @Param device_id query string false "Filter by device ID"
// @Param action query string false "Filter by action (approved, not_found, failed, search_error)"
// @Param success query boolean false "Filter by success status"
// @Param limit query integer false "Limit results (default 50)" default(50)
// @Param offset query integer false "Offset for pagination (default 0)" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /automatic/logs [get]
func (h *AutomaticHandler) GetAutomaticLogs(c *gin.Context) {
    // Получаем параметры фильтрации
    traderId := c.Query("trader_id")
    deviceId := c.Query("device_id")
    action := c.Query("action")
    successStr := c.Query("success")
    
    limitStr := c.DefaultQuery("limit", "50")
    offsetStr := c.DefaultQuery("offset", "0")
    
    // Парсим лимит и офсет
    limit := 50
    offset := 0
    
    if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 200 {
        limit = l
    }
    if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
        offset = o
    }
    
    // Исправление: правильно обрабатываем optional bool
    var success *bool
    if successStr != "" {
        if successStr == "true" {
            b := true
            success = &b
        } else if successStr == "false" {
            b := false
            success = &b
        }
    }
    
    log.Printf("📊 [LOGS] Request: trader_id=%s, device_id=%s, action=%s, success=%v, limit=%d, offset=%d",
        traderId, deviceId, action, success, limit, offset)
    
    // Конструируем фильтр - ВАЖНО: если trader_id не указан, показываем все логи
    filter := &orderpb.AutomaticLogFilter{
        DeviceId: deviceId,
        Action:   action,
        Limit:    int32(limit),
        Offset:   int32(offset),
    }
    
    // Только если указан trader_id, добавляем его в фильтр
    if traderId != "" {
        filter.TraderId = traderId
    }
    
    // Исправление: присваиваем optional bool
    if success != nil {
        filter.Success = success
    }
    
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
    defer cancel()
    
    response, err := h.orderService.GetAutomaticLogs(ctx, &orderpb.GetAutomaticLogsRequest{
        Filter: filter,
    })
    
    if err != nil {
        log.Printf("❌ [LOGS] Error fetching logs: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
        return
    }
    
    // Форматируем логи для ответа
    logs := make([]map[string]interface{}, len(response.Logs))
    for i, log := range response.Logs {
        traderID := log.TraderId
        if traderID == "" {
            traderID = "Неизвестный трейдер"
        }
        
        logs[i] = map[string]interface{}{
            "id":              log.Id,
            "device_id":       log.DeviceId,
            "trader_id":       traderID,
            "order_id":        log.OrderId,
            "amount":          log.Amount,
            "payment_system":  log.PaymentSystem,
            "direction":       log.Direction,
            "methods":         log.Methods,
            "received_at":     log.ReceivedAt.AsTime().Unix(),
            "action":          log.Action,                // approved, not_found, failed, search_error
            "success":         log.Success,
            "orders_found":    log.OrdersFound,
            "error_message":   log.ErrorMessage,
            "processing_time": log.ProcessingTime,        // миллисекунды
            "bank_name":       log.BankName,
            "card_number":     maskCardNumber(log.CardNumber),
            "created_at":      log.CreatedAt.AsTime().Unix(),
        }
    }
    
    log.Printf("✅ [LOGS] Retrieved %d logs (total: %d)", len(logs), response.Total)
    
    c.JSON(http.StatusOK, gin.H{
        "logs":   logs,
        "total":  response.Total,
        "limit":  limit,
        "offset": offset,
        "count":  len(logs),
    })
}

// ==================== DEVICE STATUS ====================
// GetDeviceStatus получает статус устройства (онлайн/оффлайн)
// @Summary Get device status
// @Description Get current online status of a device
// @Tags automatic
// @Accept json
// @Produce json
// @Param device_id query string true "Device ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /automatic/device-status [get]
func (h *AutomaticHandler) GetDeviceStatus(c *gin.Context) {
    deviceId := c.Query("device_id")
    if deviceId == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "device_id is required"})
        return
    }
    
    ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
    defer cancel()
    
    response, err := h.deviceService.GetDeviceStatus(ctx, &orderpb.GetDeviceStatusRequest{
        DeviceId: deviceId,
    })
    
    if err != nil {
        log.Printf("❌ [DEVICE_STATUS] Error: %v", err)
        if status.Code(err) == codes.NotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get device status"})
        }
        return
    }
    
    // Форматируем статус
    lastPingTime := time.Unix(response.LastPing, 0)
    onlineIndicator := "🔴 Offline"
    if response.Online {
        onlineIndicator = "🟢 Online (Last ping: " + lastPingTime.Format("15:04:05") + ")"
    }
    
    c.JSON(http.StatusOK, gin.H{
        "device_id":  response.DeviceId,
        "online":     response.Online,
        "last_ping":  response.LastPing,
        "enabled":    response.Enabled,
        "status":     map[string]interface{}{
            "online_indicator": onlineIndicator,
        },
    })
}

// GetTraderDevicesStatus получает статусы всех устройств трейдера
// @Summary Get trader devices status
// @Description Get online status of all trader's devices
// @Tags automatic
// @Accept json
// @Produce json
// @Param trader_id query string true "Trader ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /automatic/trader-devices-status [get]
func (h *AutomaticHandler) GetTraderDevicesStatus(c *gin.Context) {
    traderID := c.Query("trader_id")
    if traderID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "trader_id is required"})
        return
    }
    
    log.Printf("📱 [TRADER_DEVICES_STATUS] Request: trader_id=%s", traderID)
    
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
    defer cancel()
    
    response, err := h.deviceService.GetTraderDevicesStatus(ctx, &orderpb.GetTraderDevicesStatusRequest{
        TraderId: traderID,
    })
    
    if err != nil {
        log.Printf("❌ [TRADER_DEVICES_STATUS] Error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get devices status"})
        return
    }
    
    // Форматируем ответ
    devices := make([]map[string]interface{}, len(response.Devices))
    onlineCount := 0
    
    for i, device := range response.Devices {
        if device.Online {
            onlineCount++
        }
        
        lastPingTime := time.Unix(device.LastPing, 0)
        status := "🔴 Offline"
        if device.Online {
            status = "🟢 Online"
        }
        
        devices[i] = map[string]interface{}{
            "device_id":   device.DeviceId,
            "device_name": device.DeviceName,
            "online":      device.Online,
            "status":      status,
            "last_ping":   device.LastPing,
            "last_ping_formatted": lastPingTime.Format("2006-01-02 15:04:05"),
            "enabled":     device.Enabled,
        }
    }
    
    log.Printf("✅ [TRADER_DEVICES_STATUS] Retrieved %d devices for trader %s (%d online)", 
        len(devices), traderID, onlineCount)
    
    c.JSON(http.StatusOK, gin.H{
        "trader_id": traderID,
        "devices":   devices,
        "online_count": onlineCount,
        "total_count":  len(devices),
        "summary": map[string]interface{}{
            "all_online": onlineCount == len(devices) && len(devices) > 0,
            "any_online": onlineCount > 0,
        },
    })
}

// ==================== UTILITY FUNCTIONS ====================

// maskCardNumber маскирует номер карты для безопасности
func maskCardNumber(card string) string {
    if len(card) >= 4 {
        return "***" + card[len(card)-4:]
    }
    return card
}

// internal/handlers/automatic.go

// ==================== AUTOMATIC STATS ====================

// GetAutomaticStats получает статистику автоматической обработки
// @Summary Get automatic processing statistics
// @Description Get statistics for automatic payment processing
// @Tags automatic
// @Accept json
// @Produce json
// @Param trader_id query string true "Trader ID"
// @Param days query integer false "Number of days for statistics (default 7)" default(7)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /automatic/stats [get]
func (h *AutomaticHandler) GetAutomaticStats(c *gin.Context) {
    traderID := c.Query("trader_id")
    
    daysStr := c.DefaultQuery("days", "7")
    days, err := strconv.Atoi(daysStr)
    if err != nil || days <= 0 {
        days = 7
    }
    
    log.Printf("📊 [STATS] Request: trader_id=%s, days=%d", traderID, days)
    
    // Если trader_id не указан, возвращаем общую статистику
    if traderID == "" {
        // TODO: Реализовать общую статистику для всех трейдеров
        c.JSON(http.StatusOK, gin.H{
            "trader_id": "all",
            "period_days": days,
            "overview": map[string]interface{}{
                "total_attempts": 0,
                "successful_attempts": 0,
                "success_rate": 0,
                "approved_orders": 0,
                "not_found_count": 0,
                "failed_count": 0,
                "avg_processing_time_ms": 0,
            },
            "device_stats": map[string]interface{}{},
        })
        return
    }
    
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
    defer cancel()
    
    response, err := h.orderService.GetAutomaticStats(ctx, &orderpb.GetAutomaticStatsRequest{
        TraderId: traderID,
        Days:     int32(days),
    })
    
    if err != nil {
        log.Printf("❌ [STATS] Error fetching stats: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch statistics"})
        return
    }
    
    // Форматируем ответ
    stats := map[string]interface{}{
        "trader_id": traderID,
        "period_days": days,
        "overview": map[string]interface{}{
            "total_attempts": response.Stats.TotalAttempts,
            "successful_attempts": response.Stats.SuccessfulAttempts,
            "success_rate": calculateSuccessRate(response.Stats.TotalAttempts, response.Stats.SuccessfulAttempts),
            "approved_orders": response.Stats.ApprovedOrders,
            "not_found_count": response.Stats.NotFoundCount,
            "failed_count": response.Stats.FailedCount,
            "avg_processing_time_ms": response.Stats.AvgProcessingTime,
        },
        "device_stats": response.Stats.DeviceStats,
    }
    
    log.Printf("✅ [STATS] Retrieved stats for trader %s: %d attempts, %.1f%% success", 
        traderID, response.Stats.TotalAttempts, 
        calculateSuccessRate(response.Stats.TotalAttempts, response.Stats.SuccessfulAttempts))
    
    c.JSON(http.StatusOK, stats)
}

// ==================== RECENT ACTIVITY ====================

// GetRecentAutomaticActivity получает последние активности автоматики
// @Summary Get recent automatic activity
// @Description Get recent automatic payment processing activities
// @Tags automatic
// @Accept json
// @Produce json
// @Param trader_id query string true "Trader ID"
// @Param limit query integer false "Limit results (default 10)" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /automatic/recent-activity [get]
func (h *AutomaticHandler) GetRecentAutomaticActivity(c *gin.Context) {
    traderID := c.Query("trader_id")
    if traderID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "trader_id is required"})
        return
    }
    
    limitStr := c.DefaultQuery("limit", "10")
    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit <= 0 || limit > 50 {
        limit = 10
    }
    
    // Используем существующий endpoint логов с фильтром
    filter := &orderpb.AutomaticLogFilter{
        TraderId: traderID,
        Limit:    int32(limit),
        Offset:   0,
    }
    
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
    defer cancel()
    
    response, err := h.orderService.GetAutomaticLogs(ctx, &orderpb.GetAutomaticLogsRequest{
        Filter: filter,
    })
    
    if err != nil {
        log.Printf("❌ [RECENT_ACTIVITY] Error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recent activity"})
        return
    }
    
    // Форматируем логи для ответа
    activities := make([]map[string]interface{}, len(response.Logs))
    for i, log := range response.Logs {
        activities[i] = map[string]interface{}{
            "id":              log.Id,
            "device_id":       log.DeviceId,
            "order_id":        log.OrderId,
            "amount":          log.Amount,
            "payment_system":  log.PaymentSystem,
            "action":          log.Action,
            "success":         log.Success,
            "orders_found":    log.OrdersFound,
            "error_message":   log.ErrorMessage,
            "processing_time": log.ProcessingTime,
            "bank_name":       log.BankName,
            "received_at":     log.ReceivedAt.AsTime().Unix(),
            "created_at":      log.CreatedAt.AsTime().Unix(),
            "status_icon":     getStatusIcon(log.Success, log.Action),
        }
    }
    
    c.JSON(http.StatusOK, gin.H{
        "trader_id": traderID,
        "activities": activities,
        "count":      len(activities),
    })
}

// ==================== UTILITY FUNCTIONS ====================

// calculateSuccessRate вычисляет процент успеха
func calculateSuccessRate(total, success int64) float64 {
    if total == 0 {
        return 0
    }
    return float64(success) / float64(total) * 100
}

// getStatusIcon возвращает иконку статуса для фронтенда
func getStatusIcon(success bool, action string) string {
    if success {
        return "✅"
    }
    
    switch action {
    case "not_found":
        return "🔍"
    case "search_error":
        return "❌"
    case "failed":
        return "⚠️"
    default:
        return "❓"
    }
}

// // maskCardNumber маскирует номер карты для безопасности
// func maskCardNumber(card string) string {
//     if len(card) >= 4 {
//         return "***" + card[len(card)-4:]
//     }
//     return card
// }