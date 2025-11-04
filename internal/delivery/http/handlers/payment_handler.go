package handlers

import (
	"net/http"
	"strconv"
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
	WalletClient *client.HTTPWalletClient
	UserClient *client.UserClient
	SsoClient *client.SSOClient
}

func NewPaymentHandler(
	orderClient *client.OrderClient,
	walletClient *client.HTTPWalletClient,
	userClient *client.UserClient,
	ssoClient *client.SSOClient,
) (*PaymentHandler, error) {
	return &PaymentHandler{
		OrderClient: orderClient,
		WalletClient: walletClient,
		UserClient: userClient,
		SsoClient: ssoClient,
	}, nil
}

// @Summary sign-in handler
// @Description sign-in endpoint
// @Tags payments
// @Accept json
// @Produce json
// @Param input body merchant.LoginRequest true "user credentials"
// @Success 200 {object} merchant.LoginResponse
// @Failure 401 {object} ErrorResponse
// @Router /payments/accounts/auth/sign-in [post] 
func (h *PaymentHandler) Login(c *gin.Context) {
	var request paymentRequest.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	loginResponse, err := h.SsoClient.Login(
		request.Email,
		request.Password,
		request.TwoFaCode,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, paymentResponse.LoginResponse{
		Token: loginResponse.AccessToken,
		DateTimeExpires: loginResponse.TimeExp.AsTime().Format("2006-01-02 15:04:05"),
	})
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
	response, err := h.OrderClient.CreateOrder(&orderpb.CreateOrderRequest{
		MerchantId: payInRequest.MerchantID,
		AmountFiat: payInRequest.AmountFiat,
		Currency: payInRequest.Currency,
		ClientId: payInRequest.ClientID,
		PaymentSystem: payInRequest.PaymentSystem,
		ExpiresAt: timestamppb.New(time.Now().Add(20*time.Minute)),
		MerchantOrderId: payInRequest.MerchantOrderID,
		Shuffle: payInRequest.Shuffle,
		CallbackUrl: payInRequest.CallbackURL,
		Type: "DEPOSIT",
		BankCode: payInRequest.Issuer,
		NspkCode: "",
		Country: "",
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
			BankID: response.Order.BankDetail.BankCode,
			Bank: response.Order.BankDetail.BankCode,
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
			BankID: response.Order.BankDetail.BankCode,
			Bank: response.Order.BankDetail.BankName,
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

// @Summary Open dispute for given order
// @Description Opent dispute by order id
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "order ID"
// @Param input body paymentRequest.CreateDisputeRequest true "dispute description"
// @Success 201 {object} paymentResponse.CreateDisputeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /payments/in/h2h/{id}/arbitrage/link [post]
func (h *PaymentHandler) OpenPayInArbitrageLink(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "path param id missed"})
		return
	}
	var requestBody paymentRequest.CreateDisputeRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	disputeTtl := 30*time.Minute
	disputeID, err := h.OrderClient.CreateDispute(
		orderID,
		requestBody.ProofUrl,
		requestBody.Reason,
		disputeTtl,
		requestBody.AmountFiat,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, paymentResponse.CreateDisputeResponse{
		DisputeID: disputeID,
		ExpiresAt: time.Now().Add(disputeTtl).Local().String(),
	})
}

// @Summary Get info about dispute
// @Description Get dispute info by disputeID
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "dispute id"
// @Success 200 {object} paymentResponse.GetPayInArbitrageInfoResponse
// @Success 400 {object} ErrorResponse
// @Success 404 {object} ErrorResponse
// @Router /payments/in/h2h/{id}/arbitrage/info [get]
func (h *PaymentHandler) GetPayInArbitrageInfo(c *gin.Context) {
	disputeID := c.Param("id")
	if disputeID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "dispute id path param missed"})
		return
	}
	dispute, err := h.OrderClient.GetDisputeInfo(disputeID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, paymentResponse.GetPayInArbitrageInfoResponse{
		Dispute: paymentResponse.Dispute{
			DisputeID: dispute.DisputeID,
			OrderID: dispute.OrderID,
			ProofUrl: dispute.ProofUrl,
			DisputeReason: dispute.DisputeReason,
			DisputeStatus: dispute.DisputeStatus,
		},
	})
}

// @Summary Get balance
// @Description Get account balance
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} merchant.GetAccountBalanceResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /payments/accounts/balance [get]
func (h *PaymentHandler) GetAccountBalance(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
		return
	}

	userIDstr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID is not a string"})
		return
	}

	balance, err := h.WalletClient.GetBalance(userIDstr)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	userResp, err := h.UserClient.GetUserByID(userIDstr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get account info"})
		return
	}

	c.JSON(http.StatusOK, paymentResponse.GetAccountBalanceResponse{
		Balances: []paymentResponse.AccountBalance{
			paymentResponse.AccountBalance{
				Name: userResp.Username,
				Balance: strconv.FormatFloat(balance, 'f', 6, 64),
				Currency: "USDT",
			},
		},
	})
}

func (h *PaymentHandler) CreateRedirectPayIn(c *gin.Context) {

}

func (h *PaymentHandler) GetRedirectPayInInfo(c *gin.Context) {

}

// @Summary Get order status
// @Description Get order status
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orderId path string true "order ID"
// @Success 200 {object} merchant.GetOrderStatusResponse
// @Failure 404 {object} ErrorResponse
// @Router /payments/order/{orderId}/status [get]
func (h *PaymentHandler) GetOrderStatus(c *gin.Context) {
	orderID := c.Param("orderId")
	orderResponse, err := h.OrderClient.GetOrderByID(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to find order"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": orderResponse.Order.Status})
}

// @Summary Получить список ордеров
// @Description Получение списка ордеров с фильтрацией, сортировкой и пагинацией
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param dealId query string false "Фильтр по ID сделки"
// @Param type query string false "Тип ордера: DEPOSIT, WITHDRAWAL, PAYOUT" Enums(DEPOSIT, WITHDRAWAL, PAYOUT)
// @Param status query string false "Статус ордера" Enums(COMPLETED, CANCELED, FAILED, DISPUTE, PENDING)
// @Param timeOpeningStart query string false "Начальная дата создания (формат: 2006-01-02T15:04:05Z)"
// @Param timeOpeningEnd query string false "Конечная дата создания (формат: 2006-01-02T15:04:05Z)"
// @Param amountMin query number false "Минимальная сумма"
// @Param amountMax query number false "Максимальная сумма"
// @Param page query int false "Номер страницы (начиная с 0)" default(0)
// @Param size query int false "Размер страницы" default(10)
// @Param sort query string false "Поле сортировки: id, deal_id, time_opening, time_expires, time_complete, type, status, amount"
// @Success 200 {object} GetOrdersResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /payments/order [get]
func (h *PaymentHandler) GetOrders(c *gin.Context) {
    var params paymentRequest.GetOrdersParams
    if err := c.ShouldBindQuery(&params); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters: " + err.Error()})
        return
    }

    // Получаем merchantID из аутентификации
    merchantID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Merchant not authenticated"})
        return
    }

	merchantIDstr, ok := merchantID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid merchant ID"})
		return
	}

	userResp, err := h.UserClient.GetUserByID(merchantIDstr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to find user"})
		return
	}

    // Обрабатываем значения по умолчанию для page и size
    page := 0
    if params.Page != nil {
        page = *params.Page
    }

    size := 10 // значение по умолчанию
    if params.Size != nil {
        size = *params.Size
    }

    // Ограничиваем максимальный размер страницы
    if size > 100 {
        size = 100
    }

    // Формируем gRPC запрос
    grpcReq := &orderpb.GetOrdersRequest{
        MerchantId: merchantID.(string),
        DealId:     params.DealID,
        Type:       params.Type,
        Status:     params.Status,
        AmountMin:  params.AmountMin,
        AmountMax:  params.AmountMax,
        Sort:       params.Sort,
        Page:       int32(page),
        Size:       int32(size),
    }

    // Преобразуем временные параметры
    if params.TimeOpeningStart != nil {
        grpcReq.TimeOpeningStart = timestamppb.New(*params.TimeOpeningStart)
    }
    if params.TimeOpeningEnd != nil {
        grpcReq.TimeOpeningEnd = timestamppb.New(*params.TimeOpeningEnd)
    }
    
    grpcResp, err := h.OrderClient.GetOrders(grpcReq)
    if err != nil {
        // Обработка ошибок gRPC
        st, ok := status.FromError(err)
        if ok {
            switch st.Code() {
            case codes.InvalidArgument:
                c.JSON(http.StatusBadRequest, gin.H{"error": st.Message()})
            case codes.NotFound:
                c.JSON(http.StatusNotFound, gin.H{"error": st.Message()})
            default:
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
            }
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
        }
        return
    }

    // Преобразуем gRPC ответ в HTTP ответ
    response := toHTTPResponse(grpcResp, userResp.Login, userResp.Username)
    c.JSON(http.StatusOK, response)
}

// @Summary Withdraw crypto from account wallet
// @Description Withdraw USDT
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body merchant.WithdrawRequest true "withdraw data"
// @Success 201 {object} merchant.WithdrawResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /payments/accounts/withdraw/create [post]
func (h *PaymentHandler) Withdraw(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
		return
	}
	userIDstr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID is not a string"})
		return
	}
	var withdrawRequest paymentRequest.WithdrawRequest
	if err := c.ShouldBindJSON(&withdrawRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txHash, err := h.WalletClient.Withdraw(
		userIDstr,
		withdrawRequest.ToAddress,
		withdrawRequest.Amount,
	)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, paymentResponse.WithdrawResponse{
		TxHash: txHash,
		Amount: withdrawRequest.Amount,
		Currency: withdrawRequest.Currency,
		ToAddress: withdrawRequest.ToAddress,
	})

}