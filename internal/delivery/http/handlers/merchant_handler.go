package handlers

import (
	"net/http"
	"time"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	"github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/merchant"
	orderpb "github.com/LavaJover/shvark-order-service/proto/gen"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MerchantHandler struct {
	OrderClient *client.OrderClient
}

func NewMerchanHandler(orderClient *client.OrderClient) *MerchantHandler {
	return &MerchantHandler{
		OrderClient: orderClient,
	}
}

// @Summary Create new deposit order
// @Description Create new pay-in order
// @Tags merchant
// @Accept json
// @Produce json
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
	orderServiceRequest := orderpb.CreateOrderRequest{
		MerchantId: merchantID,
		AmountFiat: request.Amount,
		Currency: request.Currency,
		Country: "Russia",
		ClientId: "",
		ExpiresAt: timestamppb.New(time.Now().Add(20*time.Minute)),
		MerchantOrderId: request.IternalID,
		Shuffle: 0,
		CallbackUrl: request.CallbackUrl,
		BankCode: request.Issuer,
		NspkCode: request.NspkCode,
	}
	if request.IsSbp {
		orderServiceRequest.PaymentSystem = "SBP"
	}else {
		orderServiceRequest.PaymentSystem = "C2C"
	}
	orderServiceResponse, err := h.OrderClient.CreateOrder(&orderServiceRequest)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
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
		TimeExpires: orderServiceResponse.Order.ExpiresAt.Seconds,
	})
}

func (h *MerchantHandler) GetOrders(c *gin.Context) {

}

func (h *MerchantHandler) GetOrderStatus(c *gin.Context) {

}

func (h *MerchantHandler) GetAccountBalance(c *gin.Context) {

}