package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	orderpb "github.com/LavaJover/shvark-order-service/proto/gen"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AutomaticHandler struct {
	orderService *client.OrderClient
}

func NewAutomaticHandler(orderService *client.OrderClient) *AutomaticHandler {
	return &AutomaticHandler{
		orderService: orderService,
	}
}

// SMSRequest - структура для парсинга входящих уведомлений
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

func (h *AutomaticHandler) Sms(c *gin.Context) {
	var req SMSRequest
	
	if err := c.BindJSON(&req); err != nil {
		log.Printf("Ошибка парсинга JSON: %v", err)
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	// Валидация входящего уведомления
	if !h.validateSMS(req) {
		c.JSON(200, gin.H{"status": "ignored", "reason": "validation failed"})
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
		Metadata:      map[string]string{},
	}

	// Вызов order-service с retry логикой
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	response, err := h.withRetry(ctx, 3, func() (*orderpb.ProcessAutomaticPaymentResponse, error) {
		return h.orderService.ProcessAutomaticPayment(ctx, grpcReq)
	})

	if err != nil {
		log.Printf("Ошибка обработки уведомления: %v", err)
		c.JSON(500, gin.H{"error": "processing failed"})
		return
	}

	c.JSON(200, gin.H{
		"status":    "processed",
		"order_id":  response.OrderId,
		"action":    response.Action,
		"processed": response.Success,
	})
}

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

	// Проверяем направление платежа
	if req.Direction != "in" {
		return false
	}

	return true
}

func (h *AutomaticHandler) withRetry(ctx context.Context, maxRetries int, fn func() (*orderpb.ProcessAutomaticPaymentResponse, error)) (*orderpb.ProcessAutomaticPaymentResponse, error) {
	var lastErr error
	
	for i := 0; i < maxRetries; i++ {
		resp, err := fn()
		if err == nil {
			return resp, nil
		}
		
		// Проверяем, стоит ли повторять запрос
		if status.Code(err) == codes.Unavailable || status.Code(err) == codes.DeadlineExceeded {
			lastErr = err
			time.Sleep(time.Duration(i*100) * time.Millisecond)
			continue
		}
		
		return nil, err
	}
	
	return nil, lastErr
}


func (h *AutomaticHandler) Live(c *gin.Context) {
	// Используем map[string]interface{} для любого JSON
	var body map[string]interface{}
	// Парсим JSON из тела запроса
	if err := c.BindJSON(&body); err != nil {
		log.Println("Ошибка при парсинге JSON:", err)
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	// Логируем полученный реквест
	log.Println("Получен реквест:", body)

	// Можно дополнительно вывести тип каждого поля
	for k, v := range body {
		log.Printf("Поле: %s, Тип: %T, Значение: %v\n", k, v, v)
	}

	c.JSON(http.StatusOK, gin.H{
		"alive": true,
		"received": body,
	})
}

func (h *AutomaticHandler) Auth(c *gin.Context) {
	// Используем map[string]interface{} для любого JSON
	var body map[string]interface{}
	// Парсим JSON из тела запроса
	if err := c.BindJSON(&body); err != nil {
		log.Println("Ошибка при парсинге JSON:", err)
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	// Логируем полученный реквест
	log.Println("Получен реквест:", body)

	// Можно дополнительно вывести тип каждого поля
	for k, v := range body {
		log.Printf("Поле: %s, Тип: %T, Значение: %v\n", k, v, v)
	}

	// Отправляем ответ
	c.JSON(200, gin.H{
		"status": "ok",
		"received": body,
	})
}