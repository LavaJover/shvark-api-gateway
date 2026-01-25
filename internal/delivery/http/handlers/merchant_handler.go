package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/LavaJover/shvark-api-gateway/internal/client/order-service"
	ssoservice "github.com/LavaJover/shvark-api-gateway/internal/client/sso-service"
	userservice "github.com/LavaJover/shvark-api-gateway/internal/client/user-service"
	walletservice "github.com/LavaJover/shvark-api-gateway/internal/client/wallet-service"
	"github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/merchant"
	orderpb "github.com/LavaJover/shvark-order-service/proto/gen/order"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MerchantHandler struct {
	OrderClient *orderservice.OrderClient
	WalletClient *walletservice.HTTPWalletClient
	UserClient *userservice.UserClient
	SsoClient *ssoservice.SSOClient
}

func NewMerchanHandler(
	orderClient *orderservice.OrderClient,
	walletClient *walletservice.HTTPWalletClient,
	userClient *userservice.UserClient,
	ssoClient *ssoservice.SSOClient,
) *MerchantHandler {
	return &MerchantHandler{
		OrderClient: orderClient,
		WalletClient: walletClient,
		UserClient: userClient,
		SsoClient: ssoClient,
	}
}


// @Summary Create new deposit order
// @Description Create new pay-in order
// @Tags merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param accountID path string true "merchant account ID"
// @Param input body merchant.CreatePayInRequest true "new deposit order details"
// @Success 201 {object} merchant.CreatePayInResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /merchant/order/{accountID}/deposit [post]
func (h *MerchantHandler) CreatePayIn(c *gin.Context) {
	merchantID := c.Param("accountID")
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "accountID path param missed"})
		return
	}
	var request merchant.CreatePayInRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orderServiceRequest := orderpb.CreatePayInOrderRequest{
		MerchantId: merchantID,
		AmountFiat: request.Amount,
		Currency: request.Currency,
		StoreId: "2ded072b-e251-4783-9662-51fb6c713756",
		Country: "Russia",
		ClientId: "",
		ExpiresAt: timestamppb.New(time.Now().Add(20*time.Minute)),
		MerchantOrderId: request.IternalID,
		Shuffle: 0,
		CallbackUrl: request.CallbackUrl,
		BankCode: request.Issuer,
		NspkCode: request.NspkCode,
		Type: "DEPOSIT",
	}
	if request.IsSbp {
		orderServiceRequest.PaymentSystem = "SBP"
	}else {
		orderServiceRequest.PaymentSystem = "C2C"
	}
	if merchantID == "80da023c-9604-4996-afac-5d7729575622"{
		orderServiceRequest.StoreId = "2ded072b-e251-4783-9662-51fb6c713756"
	}
	if merchantID == "455854ce-27e8-405b-9cbf-83a9c9ac8109" {
		orderServiceRequest.StoreId = "b847e3f6-25d0-4c47-a68a-22393572a755"
	}
	orderServiceResponse, err := h.OrderClient.CreatePayInOrder(&orderServiceRequest)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
 	c.JSON(http.StatusCreated, merchant.CreatePayInResponse{
		OrderID: orderServiceResponse.Order.OrderId,
		CardNumber: orderServiceResponse.Order.BankDetail.CardNumber,
		PhoneNumber: orderServiceResponse.Order.BankDetail.Phone,
		HolderName: orderServiceResponse.Order.BankDetail.Owner,
		Issuer: orderServiceResponse.Order.BankDetail.BankCode,
		NspkCode: orderServiceResponse.Order.BankDetail.NspkCode,
		Amount: orderServiceResponse.Order.AmountFiat,
		AmountByCurrency: orderServiceResponse.Order.AmountCrypto,
		CurrencyRate: orderServiceResponse.Order.CryptoRubRate,
		TimeExpires: orderServiceResponse.Order.ExpiresAt.AsTime().UTC().Format("2006-01-02T15:04:05Z"),
	})
}

// @Summary Получить список ордеров
// @Description Получение списка ордеров с фильтрацией, сортировкой и пагинацией
// @Tags merchant
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
// @Router /merchant/order [get]
func (h *MerchantHandler) GetOrders(c *gin.Context) {
    var params merchant.GetOrdersParams
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

// Структуры для Swagger

// GetOrdersResponse структура успешного ответа
type GetOrdersResponse struct {
    Content          []OrderResponse `json:"content"`
    Pageable         Pageable        `json:"pageable"`
    TotalElements    int32           `json:"totalElements"`
    TotalPages       int32           `json:"totalPages"`
    Last             bool            `json:"last"`
    NumberOfElements int32           `json:"numberOfElements"`
    Size             int32           `json:"size"`
    Number           int32           `json:"number"`
    Sort             Sort            `json:"sort"`
    First            bool            `json:"first"`
    Empty            bool            `json:"empty"`
}

// OrderResponse элемент списка ордеров
type OrderResponse struct {
    ID           string     `json:"id"`
    TimeOpening  time.Time  `json:"timeOpening"`
    TimeExpires  time.Time  `json:"timeExpires"`
    TimeComplete *time.Time `json:"timeComplete,omitempty"`
    StoreName    string     `json:"storeName"`
    Type         string     `json:"type"`
    Status       string     `json:"status"`
    CurrencyRate float64    `json:"currencyRate"`
    SumInvoice   Amount     `json:"sumInvoice"`
    SumDeal      Amount     `json:"sumDeal"`
    Requisites   Requisites `json:"requisites"`
    Email        string     `json:"email"`
}

// Amount сумма денег
type Amount struct {
    Amount   float64 `json:"amount"`
    Currency string  `json:"currency"`
}

// Requisites реквизиты
type Requisites struct {
    Issuer     string `json:"issuer"`
    HolderName string `json:"holderName"`
    PhoneNumber string `json:"phoneNumber"`
	CardNumber string `json:"cardNumber"`
}

// Pageable информация о пагинации
type Pageable struct {
    Sort       Sort  `json:"sort"`
    PageNumber int32 `json:"pageNumber"`
    PageSize   int32 `json:"pageSize"`
    Offset     int32 `json:"offset"`
    Paged      bool  `json:"paged"`
    Unpaged    bool  `json:"unpaged"`
}

// Sort информация о сортировке
type Sort struct {
    Unsorted bool `json:"unsorted"`
    Sorted   bool `json:"sorted"`
    Empty    bool `json:"empty"`
}


// Преобразование gRPC ответа в HTTP структуру
func toHTTPResponse(grpcResp *orderpb.GetOrdersResponse, email, username string) *GetOrdersResponse {
    content := make([]OrderResponse, 0, len(grpcResp.Content))
    for _, item := range grpcResp.Content {
        // Преобразование времени
        var timeComplete *time.Time
        if item.TimeComplete != nil {
            tc := item.TimeComplete.AsTime()
            timeComplete = &tc
        }

        content = append(content, OrderResponse{
            ID:           item.Id,
            TimeOpening:  item.TimeOpening.AsTime(),
            TimeExpires:  item.TimeExpires.AsTime(),
            TimeComplete: timeComplete,
            StoreName:    username,
            Type:         item.Type,
            Status:       item.Status,
            CurrencyRate: item.CurrencyRate,
            SumInvoice: Amount{
                Amount:   item.SumInvoice.Amount,
                Currency: item.SumInvoice.Currency,
            },
            SumDeal: Amount{
                Amount:   item.SumDeal.Amount,
                Currency: item.SumDeal.Currency,
            },
            Requisites: Requisites{
                Issuer:     item.Requisites.Issuer,
                HolderName: item.Requisites.HolderName,
                PhoneNumber: item.Requisites.PhoneNumber,
				CardNumber: item.Requisites.CardNumber,
            },
            Email: email,
        })
    }

    return &GetOrdersResponse{
        Content:          content,
        Pageable:         toHTTPPageable(grpcResp.Pageable),
        TotalElements:    grpcResp.TotalElements,
        TotalPages:       grpcResp.TotalPages,
        Last:             grpcResp.Last,
        NumberOfElements: grpcResp.NumberOfElements,
        Size:             grpcResp.Size,
        Number:           grpcResp.Number,
        Sort:             toHTTPSort(grpcResp.Sort),
        First:            grpcResp.First,
        Empty:            grpcResp.Empty,
    }
}

func toHTTPPageable(p *orderpb.Pageable) Pageable {
    return Pageable{
        Sort:       toHTTPSort(p.Sort),
        PageNumber: p.PageNumber,
        PageSize:   p.PageSize,
        Offset:     p.Offset,
        Paged:      p.Paged,
        Unpaged:    p.Unpaged,
    }
}

func toHTTPSort(s *orderpb.Sort) Sort {
    return Sort{
        Unsorted: s.Unsorted,
        Sorted:   s.Sorted,
        Empty:    s.Empty,
    }
}

// @Summary Get order status
// @Description Get order status
// @Tags merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param iternalId path string true "merchant order ID"
// @Success 200 {object} merchant.GetOrderStatusResponse
// @Failure 404 {object} ErrorResponse
// @Router /merchant/order/{iternalId}/status [get]
func (h *MerchantHandler) GetOrderStatus(c *gin.Context) {
	iternalID := c.Param("iternalId")
	orderResponse, err := h.OrderClient.GetOrderByMerchantOrderID(iternalID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to find order"})
		return
	}
	c.JSON(http.StatusOK, merchant.GetOrderStatusResponse{
		Status: orderResponse.Order.Status,
	})
}

// @Summary Get balance
// @Description Get account balance
// @Tags merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} merchant.GetAccountBalanceResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /merchant/accounts/balance [get]
func (h *MerchantHandler) GetAccountBalance(c *gin.Context) {
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

	c.JSON(http.StatusOK, merchant.GetAccountBalanceResponse{
		Balances: []merchant.AccountBalance{
			merchant.AccountBalance{
				Name: userResp.Username,
				Balance: strconv.FormatFloat(balance, 'f', 6, 64),
				Currency: "USDT",
			},
		},
	})

}

// @Summary Withdraw crypto from account wallet
// @Description Withdraw USDT
// @Tags merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body merchant.WithdrawRequest true "withdraw data"
// @Success 201 {object} merchant.WithdrawResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /merchant/accounts/withdraw/create [post]
func (h *MerchantHandler) Withdraw(c *gin.Context) {
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
	var withdrawRequest merchant.WithdrawRequest
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

	c.JSON(http.StatusCreated, merchant.WithdrawResponse{
		TxHash: txHash,
		Amount: withdrawRequest.Amount,
		Currency: withdrawRequest.Currency,
		ToAddress: withdrawRequest.ToAddress,
	})

}

// @Summary sign-in handler
// @Description sign-in endpoint
// @Tags merchant
// @Accept json
// @Produce json
// @Param input body merchant.LoginRequest true "user credentials"
// @Success 200 {object} merchant.LoginResponse
// @Failure 401 {object} ErrorResponse
// @Router /merchant/auth/sign-in [post] 
func (h *MerchantHandler) Login(c *gin.Context) {
	var request merchant.LoginRequest

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
	c.JSON(http.StatusOK, merchant.LoginResponse{
		Token: loginResponse.AccessToken,
		DateTimeExpires: loginResponse.TimeExp.AsTime().Format("2006-01-02 15:04:05"),
	})
}


// @Summary Get banks
// @Description Get available banks
// @Tags merchant
// @Accept json
// @Produce json
// @Success 200
// @Router /merchant/banks [get]
func (h *MerchantHandler) GetBanks(c *gin.Context) {
	banks := []merchant.Bank{
		// {
		// 	Code: "MTS-SIM",
		// 	Name: "МТС SIM",
		// 	NspkCode: "",
		// },
		// {
		// 	Code: "BEELINE-SIM",
		// 	Name: "Билайн SIM",
		// 	NspkCode: "",
		// },
		{
			Code: "sberbank",
			Name: "Сбербанк",
			NspkCode: "100000000111",
		  },
		  {
			Code: "tinkoff",
			Name: "Т-Банк",
			NspkCode: "100000000004",
		  },
		  {
			Code: "vtb",
			Name: "Банк ВТБ",
			NspkCode: "110000000005",
		  },
		  {
			Code: "alfabank",
			Name: "АЛЬФА-БАНК",
			NspkCode: "100000000008",
		  },
		  {
			Code: "raiffeisenbank",
			Name: "Райффайзенбанк",
			NspkCode: "100000000007",
		  },
		  {
			Code: "bank_open",
			Name: "Банк ОТКРЫТИЕ",
			NspkCode: "100000000015",
		  },
		  {
			Code: "gazprombank",
			Name: "Газпромбанк",
			NspkCode: "100000000001",
		  },
		  {
			Code: "promsvyaz",
			Name: "Промсвязьбанк",
			NspkCode: "100000000010",
		  },
		  {
			Code: "sovkom",
			Name: "Совкомбанк",
			NspkCode: "100000000013",
		  },
		  {
			Code: "ros_bank",
			Name: "РОСБАНК",
			NspkCode: "100000000012",
		  },
		  {
			Code: "rsb",
			Name: "Банк Русский Стандарт",
			NspkCode: "100000000014",
		  },
		  {
			Code: "absolute_bank",
			Name: "АКБ Абсолют Банк",
			NspkCode: "100000000047",
		  },
		  {
			Code: "home_bank",
			Name: "Хоум кредит",
			NspkCode: "100000000024",
		  },
		  {
			Code: "otp_bank",
			Name: "ОТП Банк",
			NspkCode: "100000000018",
		  },
		  {
			Code: "uralsib",
			Name: "БАНК УРАЛСИБ",
			NspkCode: "100000000026",
		  },
		  {
			Code: "ak_bars_bank",
			Name: "АК БАРС БАНК",
			NspkCode: "100000000006",
		  },
		  {
			Code: "fora",
			Name: "АКБ ФОРА-БАНК",
			NspkCode: "100000000217",
		  },
		  {
			Code: "rost_finance",
			Name: "КБ РостФинанс",
			NspkCode: "100000000098",
		  },
		  {
			Code: "ozon",
			Name: "Озон Банк (Ozon)",
			NspkCode: "100000000273",
		  },
		  {
			Code: "unistream",
			Name: "КБ ЮНИСТРИМ",
			NspkCode: "100000000042",
		  },
		  {
			Code: "mts",
			Name: "МТС-Банк",
			NspkCode: "100000000017",
		  },
		  {
			Code: "tkb",
			Name: "ТрансКапиталБанк",
			NspkCode: "100000000034",
		  },
		  {
			Code: "pochta",
			Name: "Почта Банк",
			NspkCode: "100000000016",
		  },
		  {
			Code: "rncb",
			Name: "РНКБ Банк",
			NspkCode: "100000000011",
		  },
		  {
			Code: "genbank",
			Name: "ГЕНБАНК",
			NspkCode: "100000000037",
		  },
		  {
			Code: "cifra",
			Name: "Цифра банк",
			NspkCode: "100000000265",
		  },
		  {
			Code: "ingo",
			Name: "Ингосстрах Банк",
			NspkCode: "100000000078",
		  },
		  {
			Code: "svoi",
			Name: "Свой Банк",
			NspkCode: "100000000006",
		  },
		  {
			Code: "avangard",
			Name: "АКБ АВАНГАРД",
			NspkCode: "100000000028",
		  },
		  {
			Code: "rencredit",
			Name: "КБ Ренессанс Кредит",
			NspkCode: "100000000032",
		  },
		  {
			Code: "solid",
			Name: "КБ Солидарность",
			NspkCode: "100000000121",
		  },
		  {
			Code: "pucbr",
			Name: "ПУ Банк России",
			NspkCode: "100000000027",
		  },
		  {
			Code: "expobank",
			Name: "Экспобанк",
			NspkCode: "100000000044",
		  },
		  {
			Code: "apkbank",
			Name: "КБ АГРОПРОМКРЕДИТ",
			NspkCode: "100000000118",
		  },
		  {
			Code: "bktb",
			Name: "Кубаньторгбанк",
			NspkCode: "100000000180",
		  },
		  {
			Code: "bystrobank",
			Name: "БыстроБанк",
			NspkCode: "100000000092",
		  },
		  {
			Code: "nico-bank",
			Name: "НИКО-БАНК",
			NspkCode: "100000000115",
		  },
		  {
			Code: "okbank",
			Name: "Банк Объединенный капитал",
			NspkCode: "100000000182",
		  },
		  {
			Code: "pscb",
			Name: "Банк ПСКБ",
			NspkCode: "100000000087",
		  },
		  {
			Code: "zarech",
			Name: "Банк Заречье",
			NspkCode: "100000000205",
		  },
		  {
			Code: "zemsky",
			Name: "Земский банк",
			NspkCode: "100000000066",
		  },
		  {
			Code: "abr",
			Name: "АБ РОССИЯ",
			NspkCode: "100000000095",
		  },
		  {
			Code: "bksbank",
			Name: "БКС Банк",
			NspkCode: "100000000041",
		  },
		  {
			Code: "chelinvest",
			Name: "ЧЕЛЯБИНВЕСТБАНК",
			NspkCode: "100000000094",
		  },
		  {
			Code: "databank",
			Name: "Датабанк",
			NspkCode: "100000000070",
		  },
		  {
			Code: "domrfbank",
			Name: "Банк ДОМ.РФ",
			NspkCode: "100000000082",
		  },
		  {
			Code: "energobank",
			Name: "АКБ Энергобанк",
			NspkCode: "100000000159",
		  },
		  {
			Code: "forshtadt",
			Name: "АКБ Форштадт",
			NspkCode: "100000000081",
		  },
		  {
			Code: "gaztransbank",
			Name: "Газтрансбанк",
			NspkCode: "100000000183",
		  },
		  {
			Code: "gebank",
			Name: "Газэнергобанк",
			NspkCode: "100000000043",
		  },
		  {
			Code: "in-bank",
			Name: "Инбанк",
			NspkCode: "100000000196",
		  },
		  {
			Code: "iturupbank",
			Name: "Банк ИТУРУП",
			NspkCode: "100000000158",
		  },
		  {
			Code: "kbb",
			Name: "Кузнецкбизнесбанк",
			NspkCode: "100000000195",
		  },
		  {
			Code: "kbhmb",
			Name: "Хакасский муниципальный банк",
			NspkCode: "100000000127",
		  },
		  {
			Code: "kkbank",
			Name: "КБ Кубань Кредит",
			NspkCode: "100000000050",
		  },
		  {
			Code: "koshelev-bank",
			Name: "КОШЕЛЕВ-БАНК",
			NspkCode: "100000000146",
		  },
		  {
			Code: "kremlinbank",
			Name: "Банк Кремлевский",
			NspkCode: "100000000201",
		  },
		  {
			Code: "maritimebank",
			Name: "МОРСКОЙ БАНК",
			NspkCode: "100000000171",
		  },
		  {
			Code: "mcbank",
			Name: "БАНК МОСКВА-СИТИ",
			NspkCode: "100000000234",
		  },
		  {
			Code: "metallinvestbank",
			Name: "Металлинвестбанк",
			NspkCode: "100000000046",
		  },
		  {
			Code: "akibank",
			Name: "АКИБАНК",
			NspkCode: "100000000107",
		  },
		  {
			Code: "modulbank",
			Name: "КБ Модульбанк",
			NspkCode: "100000000099",
		  },
		  {
			Code: "mp-bank",
			Name: "МП Банк",
			NspkCode: "100000000169",
		  },
		  {
			Code: "nkbank",
			Name: "НК Банк",
			NspkCode: "100000000233",
		  },
		  {
			Code: "norvikbank",
			Name: "Норвик Банк",
			NspkCode: "100000000202",
		  },
		  {
			Code: "novikom",
			Name: "Банк НОВИКОМ (НОВИКОМБАНК)",
			NspkCode: "100000000177",
		  },
		  {
			Code: "novobank",
			Name: "УКБ Новобанк",
			NspkCode: "100000000222",
		  },
		  {
			Code: "nsbank",
			Name: "НС Банк",
			NspkCode: "100000000071",
		  },
		  {
			Code: "orbank",
			Name: "БАНК ОРЕНБУРГ",
			NspkCode: "100000000124",
		  },
		  {
			Code: "crediteurope",
			Name: "Кредит Европа Банк (Россия)",
			NspkCode: "100000000027",
		  },
		  {
			Code: "pskb",
			Name: "СКБ Приморья Примсоцбанк",
			NspkCode: "100000000088",
		  },
		  {
			Code: "realistbank",
			Name: "РЕАЛИСТ БАНК",
			NspkCode: "100000000232",
		  },
		  {
			Code: "resocreditbank",
			Name: "Банк РЕСО Кредит",
			NspkCode: "100000000187",
		  },
		  {
			Code: "sinko-bank",
			Name: "КБ СИНКО-БАНК",
			NspkCode: "100000000148",
		  },
		  {
			Code: "socium-bank",
			Name: "СОЦИУМ-БАНК",
			NspkCode: "100000000223",
		  },
		  {
			Code: "tatsotsbank",
			Name: "ТАТСОЦБАНК",
			NspkCode: "100000000189",
		  },
		  {
			Code: "timerbank",
			Name: "Тимер Банк",
			NspkCode: "100000000144",
		  },
		  {
			Code: "transstroybank",
			Name: "АКБ Трансстройбанк",
			NspkCode: "100000000197",
		  },
		  {
			Code: "agros",
			Name: "Банк Агророс",
			NspkCode: "100000000102",
		  },
		  {
			Code: "akcept",
			Name: "Банк Акцепт",
			NspkCode: "100000000135",
		  },
		  {
			Code: "alefbank",
			Name: "АКБ Алеф-Банк",
			NspkCode: "100000000113",
		  },
		  {
			Code: "aresbank",
			Name: "КБ АРЕСБАНК",
			NspkCode: "100000000129",
		  },
		  {
			Code: "bancaintesa",
			Name: "Банк Интеза",
			NspkCode: "100000000170",
		  },
		  {
			Code: "bank-hlynov",
			Name: "КБ Хлынов",
			NspkCode: "100000000056",
		  },
		  {
			Code: "bankofkazan",
			Name: "КБЭР Банк Казани",
			NspkCode: "100000000191",
		  },
		  {
			Code: "bspb",
			Name: "Банк Санкт-Петербург",
			NspkCode: "100000000029",
		  },
		  {
			Code: "centrinvest",
			Name: "КБ Центр-инвест",
			NspkCode: "100000000059",
		  },
		  {
			Code: "cfb",
			Name: "Банк БКФ",
			NspkCode: "100000000227",
		  },
		  {
			Code: "chelindbank",
			Name: "ЧЕЛИНДБАНК",
			NspkCode: "100000000106",
		  },
		  {
			Code: "coalmetbank",
			Name: "Углеметбанк",
			NspkCode: "100000000093",
		  },
		  {
			Code: "creditural",
			Name: "Кредит Урал Банк",
			NspkCode: "100000000064",
		  },
		  {
			Code: "dcapital",
			Name: "Банк Развитие-Столица",
			NspkCode: "100000000172",
		  },
		  {
			Code: "dvbank",
			Name: "Дальневосточный банк",
			NspkCode: "100000000083",
		  },
		  {
			Code: "el-plat",
			Name: "ПНКО ЭЛПЛАТ",
			NspkCode: "100000000086",
		  },
		  {
			Code: "energotransbank",
			Name: "КБ ЭНЕРГОТРАНСБАНК",
			NspkCode: "100000000139",
		  },
		  {
			Code: "finam",
			Name: "Банк ФИНАМ",
			NspkCode: "100000000040",
		  },
		  {
			Code: "gibank",
			Name: "КБ Гарант-Инвест",
			NspkCode: "100000000112",
		  },
		  {
			Code: "gorbank",
			Name: "ГОРБАНК",
			NspkCode: "100000000125",
		  },
		  {
			Code: "gutabank",
			Name: "ГУТА-БАНК",
			NspkCode: "100000000149",
		  },
		  {
			Code: "icbru",
			Name: "ИК Банк",
			NspkCode: "100000000122",
		  },
		  {
			Code: "lanta",
			Name: "АКБ Ланта-Банк",
			NspkCode: "100000000245",
		  },
		  {
			Code: "lockobank",
			Name: "КБ ЛОКО-Банк",
			NspkCode: "100000000161",
		  },
		  {
			Code: "mcbankrus",
			Name: "МС Банк Рус",
			NspkCode: "100000000229",
		  },
		  {
			Code: "metcom",
			Name: "МЕТКОМБАНК",
			NspkCode: "100000000136",
		  },
		  {
			Code: "mspbank",
			Name: "МС Примбанк",
			NspkCode: "100000000255",
		  },
		  {
			Code: "nipbank",
			Name: "Нацинвестпромбанк",
			NspkCode: "100000000185",
		  },
		  {
			Code: "ns-bank",
			Name: "Банк Национальный стандарт",
			NspkCode: "100000000243",
		  },
		  {
			Code: "nskbl",
			Name: "Банк Левобережный",
			NspkCode: "100000000052",
		  },
		  {
			Code: "poidem",
			Name: "КБ Пойдём",
			NspkCode: "100000000103",
		  },
		  {
			Code: "primbank",
			Name: "АКБ Приморье",
			NspkCode: "100000000226",
		  },
		  {
			Code: "probank",
			Name: "ПроБанк",
			NspkCode: "100000000117",
		  },
		  {
			Code: "rdb",
			Name: "РосДорБанк",
			NspkCode: "100000000084",
		  },
		  {
			Code: "sdm",
			Name: "СДМ-Банк",
			NspkCode: "100000000069",
		  },
		  {
			Code: "sevnb",
			Name: "Северный Народный Банк",
			NspkCode: "100000000208",
		  },
		  {
			Code: "sibsoc",
			Name: "СИБСОЦБАНК",
			NspkCode: "100000000166",
		  },
		  {
			Code: "sngb",
			Name: "БАНК СНГБ",
			NspkCode: "100000000091",
		  },
		  {
			Code: "tavrich",
			Name: "Таврический Банк",
			NspkCode: "100000000173",
		  },
		  {
			Code: "tenderbank",
			Name: "АКБ ТЕНДЕР-БАНК",
			NspkCode: "100000000175",
		  },
		  {
			Code: "thbank",
			Name: "Тольяттихимбанк",
			NspkCode: "100000000152",
		  },
		  {
			Code: "tpsbank",
			Name: "Томскпромстройбанк",
			NspkCode: "100000000206",
		  },
		  {
			Code: "unicreditbank",
			Name: "ЮниКредит Банк",
			NspkCode: "100000000030",
		  },
		  {
			Code: "uralfd",
			Name: "КБ Урал ФД",
			NspkCode: "100000000151",
		  },
		  {
			Code: "vbrr",
			Name: "Банк ВБРР",
			NspkCode: "100000000049",
		  },
		  {
			Code: "venets-bank",
			Name: "Банк Венец",
			NspkCode: "100000000153",
		  },
		  {
			Code: "vfbank",
			Name: "КБ ВНЕШФИНБАНК",
			NspkCode: "100000000248",
		  },
		  {
			Code: "zenit",
			Name: "Банк ЗЕНИТ",
			NspkCode: "100000000045",
		  },
		  {
			Code: "yoomoney",
			Name: "НКО ЮМани",
			NspkCode: "100000000022",
		  },
		  {
			Code: "avtofinbank",
			Name: "Авто Финанс Банк",
			NspkCode: "100000000253",
		  },
		  {
			Code: "avtotorgbank",
			Name: "Автоторгбанк",
			NspkCode: "100000000181",
		  },
		  {
			Code: "aikb-enisejskij-obedinennyj-bank",
			Name: "АИКБ Енисейский объединенный банк",
			NspkCode: "100000000258",
		  },
		  {
			Code: "bank-sinara",
			Name: "Банк Синара",
			NspkCode: "100000000003",
		  },
		  {
			Code: "bank-uralfinans",
			Name: "Уралфинанс",
			NspkCode: "100000000096",
		  },
		  {
			Code: "bank-centrokredit",
			Name: "ЦентроКредит",
			NspkCode: "100000000231",
		  },
		  {
			Code: "is-bank",
			Name: "ИС Банк",
			NspkCode: "100000000239",
		  },
		  {
			Code: "kb-dolinsk",
			Name: "КБ Долинск",
			NspkCode: "100000000270",
		  },
		  {
			Code: "novyj-vek",
			Name: "КБ Новый век",
			NspkCode: "100000000067",
		  },
		  {
			Code: "ukb-belgorodsocbank",
			Name: "УКБ Белгородсоцбанк",
			NspkCode: "100000000225",
		  },
		  {
			Code: "avtogradbank",
			Name: "Автоградбанк",
			NspkCode: "100000000130",
		  },
		  {
			Code: "bank-ekaterinburg",
			Name: "Банк Екатеринбург",
			NspkCode: "100000000090",
		  },
		  {
			Code: "bank-sgb",
			Name: "БАНК СГБ",
			NspkCode: "100000000219",
		  },
		  {
			Code: "bank-siab",
			Name: "Банк СИАБ",
			NspkCode: "100000000278",
		  },
		  {
			Code: "dzhej-jend-ti-bank",
			Name: "Джей энд Ти Банк (АО)",
			NspkCode: "100000000213",
		  },
		  {
			Code: "mb-bank",
			Name: "МБ Банк",
			NspkCode: "100000000140",
		  },
		  {
			Code: "smp-bank",
			Name: "СМП Банк",
			NspkCode: "100000000036",
		  },
		  {
			Code: "tochka-otkrytie",
			Name: "ТОЧКА (ФК Открытие)",
			NspkCode: "100000000284",
		  },
		  {
			Code: "bsdbank",
			Name: "Черноморский банк развития",
			NspkCode: "100000000215",
		  },
		  {
			Code: "klookva",
			Name: "Клюква",
			NspkCode: "100000000154",
		  },
		  {
			Code: "blanc",
			Name: "Бланк банк",
			NspkCode: "100000000053",
		  },
		  {
			Code: "tkbbank",
			Name: "ТКБ БАНК",
			NspkCode: "100000000034",
		  },
		  {
			Code: "severgazbank",
			Name: "Севергазбанк",
			NspkCode: "100000000219",
		  },
		  {
			Code: "nrb",
			Name: "АКБ НРБанк",
			NspkCode: "100000000184",
		  },
		  {
			Code: "finstarbank",
			Name: "ФИНСТАР БАНК",
			NspkCode: "100000000278",
		  },
		  {
			Code: "round",
			Name: "банк Раунд",
			NspkCode: "100000000247",
		  },
		  {
			Code: "dtb1",
			Name: "Первый Дортрансбанк",
			NspkCode: "100000000137",
		  },
		  {
			Code: "cmrbank",
			Name: "ЦМРБанк",
			NspkCode: "100000000282",
		  },
		  {
			Code: "plait",
			Name: "Плайт",
			NspkCode: "100000000296",
		  },
		  {
			Code: "bankorange",
			Name: "Банк Оранжевый",
			NspkCode: "100000000286",
		  },
		  {
			Code: "yarinterbank",
			Name: "ИКБР ЯРИНТЕРБАНК",
			NspkCode: "100000000293",
		  },
		  {
			Code: "yandexbank",
			Name: "Яндекс Банк",
			NspkCode: "100000000150",
		  },
		  {
			Code: "bank-mba-moskva",
			Name: "Банк МБА МОСКВА",
			NspkCode: "100000000192",
		  },
		  {
			Code: "bank-ipb",
			Name: "Банк ИПБ",
			NspkCode: "100000000236",
		  },
		  {
			Code: "bank-jelita",
			Name: "банк Элита",
			NspkCode: "100000000266",
		  },
		  {
			Code: "bank-vologzhanin",
			Name: "Банк Вологжанин",
			NspkCode: "100000000257",
		  },
		  {
			Code: "drajv-klik-bank",
			Name: "Драйв Клик Банк",
			NspkCode: "100000000250",
		  },
		  {
			Code: "rosselhozbank",
			Name: "Россельхозбанк",
			NspkCode: "100000000020",
		  },
		  {
			Code: "hajs",
			Name: "Хайс",
			NspkCode: "100000000272",
		  },
		  {
			Code: "jes-bi-aj-bank",
			Name: "Эс-Би-Ай Банк",
			NspkCode: "100000000105",
		  },
		  {
			Code: "rusnarbank",
			Name: "РУСНАРБАНК",
			NspkCode: "100000000194",
		  },
		  {
			Code: "bank-saratov",
			Name: "Банк Саратов",
			NspkCode: "100000000126",
		  },
		  {
			Code: "pervyj-investicionnyj-bank",
			Name: "Первый Инвестиционный Банк",
			NspkCode: "100000000174",
		  },
		  {
			Code: "akb-derzhava",
			Name: "АКБ Держава",
			NspkCode: "100000000235",
		  },
		  {
			Code: "kb-strojlesbank",
			Name: "КБ Стройлесбанк",
			NspkCode: "100000000193",
		  },
		  {
			Code: "tojota-bank",
			Name: "Тойота Банк",
			NspkCode: "100000000138",
		  },
		  {
			Code: "kb-moskommercbank",
			Name: "КБ Москоммерцбанк",
			NspkCode: "100000000110",
		  },
		  {
			Code: "uralprombank",
			Name: "УРАЛПРОМБАНК",
			NspkCode: "100000000142",
		  },
		  {
			Code: "sitibank",
			Name: "Ситибанк",
			NspkCode: "100000000128",
		  },
		  {
			Code: "bank-aleksandrovskij",
			Name: "Банк АЛЕКСАНДРОВСКИЙ",
			NspkCode: "100000000211",
		  },
		  {
			Code: "mezhdunarodnyj-finansovyj-klub",
			Name: "МЕЖДУНАРОДНЫЙ ФИНАНСОВЫЙ КЛУБ",
			NspkCode: "100000000203",
		  },
		  {
			Code: "mkb",
			Name: "Московский кредитный банк",
			NspkCode: "100000000025",
		  },
		  {
			Code: "nbd-bank",
			Name: "НБД-Банк",
			NspkCode: "100000000134",
		  },
		  {
			Code: "jandeks-bank",
			Name: "Яндекс Банк",
			NspkCode: "100000000150",
		  },
		  {
			Code: "nokssbank",
			Name: "НОКССБАНК",
			NspkCode: "100000000062",
		  },
		  {
			Code: "vuz-bank",
			Name: "ВУЗ-банк",
			NspkCode: "100000000215",
		  },
		  {
			Code: "bank-bzhf",
			Name: "Банк БЖФ",
			NspkCode: "100000000260",
		  },
		  {
			Code: "jug-investbank",
			Name: "ЮГ-Инвестбанк",
			NspkCode: "100000000160",
		  },
		  {
			Code: "kb-krokus-bank",
			Name: "КБ Крокус Банк",
			NspkCode: "100000000212",
		  },
		  {
			Code: "vladbiznesbank",
			Name: "ВЛАДБИЗНЕСБАНК",
			NspkCode: "100000000058",
		  },
		  {
			Code: "bank-avers",
			Name: "Банк Аверс",
			NspkCode: "100000000154",
		  },
		  {
			Code: "wbbank",
			Name: "Вайлдберриз Банк",
			NspkCode: "100000000259",
		  },
		  {
			Code: "ishbank",
			Name: "ИШБАНК",
			NspkCode: "100000000199",
		  },
		  {
			Code: "almazjergijenbank",
			Name: "Алмазэргиэнбанк",
			NspkCode: "",
		  },
		  {
			Code: "aziatsko-tihookeanskij-bank",
			Name: "Азиатско-Тихоокеанский Банк",
			NspkCode: "100000000108",
		  },
		  {
			Code: "agroros",
			Name: "Банк Агророс",
			NspkCode: "100000000102",
		  },
		  {
			Code: "bbr-bank",
			Name: "ББР Банк",
			NspkCode: "100000000133",
		  },
		  {
			Code: "solid-bank",
			Name: "Солид Банк",
			NspkCode: "100000000230",
		  },
		  {
			Code: "rus-universalbank",
			Name: "Русьуниверсалбанк",
			NspkCode: "100000000165",
		  },
		  {
			Code: "akb-slavija",
			Name: "АКБ СЛАВИЯ",
			NspkCode: "100000000200",
		  },
		  {
			Code: "akb-evrofinans-mosnarbank",
			Name: "АКБ ЕВРОФИНАНС МОСНАРБАНК",
			NspkCode: "100000000167",
		  },
		  {
			Code: "ubrib",
			Name: "Банк УБРиР",
			NspkCode: "100000000031",
		  },
		  {
			Code: "bank-raund",
			Name: "банк Раунд",
			NspkCode: "100000000247",
		  },
		  {
			Code: "tochka-bank",
			Name: "Точка Банк",
			NspkCode: "100000000284",
		  },
		  {
			Code: "prio-vneshtorgbank",
			Name: "Прио-Внешторгбанк",
			NspkCode: "100000000228",
		  },
		  {
			Code: "bank-snezhinskij",
			Name: "Банк Снежинский",
			NspkCode: "100000000163",
		  },
		  {
			Code: "moskombank",
			Name: "МОСКОМБАНК",
			NspkCode: "100000000176",
		  },
	}
	c.JSON(http.StatusOK, banks)
}