package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	walletRequest "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/wallet/request"
	walletResponse "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/wallet/response"
	"github.com/gin-gonic/gin"
)

type WalletHandler struct {

}

func NewWalletHandler() (*WalletHandler, error) {
	return &WalletHandler{}, nil
}

// @Summary Create new wallet
// @Description Create new wallet
// @Tags wallets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body walletRequest.CreateWalletRequest true "New wallet data"
// @Success 201 {object} walletResponse.CreateWalletErrorResponse
// @Failure 400 {object} walletResponse.CreateWalletErrorResponse
// @Failure 500 {object} walletResponse.CreateWalletErrorResponse
// @Router /wallets/create [post]
func (h *WalletHandler) CreateWallet(c *gin.Context) {
	var request walletRequest.CreateWalletRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	proxyRequestBody, err := json.Marshal(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "json marshal failed"})
		return
	}

	proxyResp, err := http.Post("http://localhost:3000/wallets/create", "application/json", bytes.NewBuffer(proxyRequestBody))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "wallet-service unavailable"})
		return
	}
	defer proxyResp.Body.Close()

	proxyRespBody, err := io.ReadAll(proxyResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read wallet-service reponse body"})
		return
	}

	var response walletResponse.CreateWalletResponse
	if err = json.Unmarshal(proxyRespBody, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "json unmarshal failed"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// @Summary Freeze crypto
// @Description Freeze crypto
// @Tags wallets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param imput body walletRequest.FreezeRequest true "wallet data"
// @Success 200 {object} walletResponse.FreezeResponse
// @Failure 400 {object} walletResponse.FreezeErrorResponse
// @Failure 500 {object} walletResponse.FreezeErrorResponse
// @Failure 502 {object} walletResponse.FreezeErrorResponse
// @Router /wallets/freeze [post]
func (h *WalletHandler) Freeze(c *gin.Context) {
	var request walletRequest.FreezeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	proxyRequestBody, err := json.Marshal(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "json marshal failed"})
		return
	}

	proxyResp, err := http.Post("http://localhost:3000/wallets/freeze", "application/json", bytes.NewBuffer(proxyRequestBody))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "wallet-service unavailable"})
		return
	}
	defer proxyResp.Body.Close()

	proxyRespBody, err := io.ReadAll(proxyResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read wallet-service reponse body"})
		return
	}

	var response walletResponse.FreezeResponse
	if err = json.Unmarshal(proxyRespBody, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "json unmarshal failed"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Release crypto
// @Description Release crypto
// @Tags wallets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body walletRequest.ReleaseRequest true "wallet data"
// @Success 200 {object} walletResponse.ReleaseResponse
// @Failure 400 {object} walletResponse.ReleaseErrorResponse
// @Failure 500 {object} walletResponse.ReleaseErrorResponse
// @Failure 502 {object} walletResponse.ReleaseErrorResponse
// @Router /wallets/release [post]
func (h *WalletHandler) Release(c *gin.Context) {
	var request walletRequest.ReleaseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	proxyRequestBody, err := json.Marshal(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "json marshal failed"})
		return
	}

	proxyResp, err := http.Post("http://localhost:3000/wallets/release", "application/json", bytes.NewBuffer(proxyRequestBody))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "wallet-service unavailable"})
		return
	}
	defer proxyResp.Body.Close()

	proxyRespBody, err := io.ReadAll(proxyResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read wallet-service reponse body"})
		return
	}

	var response walletResponse.ReleaseResponse
	if err = json.Unmarshal(proxyRespBody, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "json unmarshal failed"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Withdraw crypto
// @Description Withdraw crypto
// @Tags wallets
// @Security BearerAuth
// @Param input body walletRequest.WithdrawRequest true "wallet data"
// @Accept json
// @Produce json
// @Success 200 {object} walletResponse.WithdrawResponse
// @Failure 400 {object} walletResponse.WithdrawErrorResponse
// @Failure 500 {object} walletResponse.WithdrawErrorResponse
// @Failure 502 {object} walletResponse.WithdrawErrorResponse
// @Router /wallets/withdraw [post]
func (h *WalletHandler) Withdraw(c *gin.Context) {
	var request walletRequest.WithdrawRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	proxyRequestBody, err := json.Marshal(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "json marshal failed"})
		return
	}

	proxyResp, err := http.Post("http://localhost:3000/wallets/withdraw", "application/json", bytes.NewBuffer(proxyRequestBody))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "wallet-service unavailable"})
		return
	}
	defer proxyResp.Body.Close()

	proxyRespBody, err := io.ReadAll(proxyResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read wallet-service reponse body"})
		return
	}

	var response walletResponse.WithdrawResponse
	if err = json.Unmarshal(proxyRespBody, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "json unmarshal failed"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Deposit crypto off-chain
// @Description Deposit crypto off chain
// @Tags wallets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body walletRequest.DepositRequest true "wallet data"
// @Success 200 {object} walletResponse.DepositResponse
// @Failure 400 {object} walletResponse.DepositErrorResponse
// @Failure 500 {object} walletResponse.DepositErrorResponse
// @Failure 502 {object} walletResponse.DepositErrorResponse
// @Router /wallets/deposit [post]
func (h *WalletHandler) Deposit(c *gin.Context) {
		var request walletRequest.DepositRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	proxyRequestBody, err := json.Marshal(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "json marshal failed"})
		return
	}

	proxyResp, err := http.Post("http://localhost:3000/wallets/deposit", "application/json", bytes.NewBuffer(proxyRequestBody))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "wallet-service unavailable"})
		return
	}
	defer proxyResp.Body.Close()

	proxyRespBody, err := io.ReadAll(proxyResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read wallet-service reponse body"})
		return
	}

	var response walletResponse.DepositResponse
	if err = json.Unmarshal(proxyRespBody, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "json unmarshal failed"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get trader transactions history
// @Description Get trader transaction history
// @Tags wallets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param traderID path string true "TraderID"
// @Success 200 {object} walletResponse.GetTraderHistoryResponse
// @Failure 400 {object} walletResponse.GetTraderHistoryErrorResponse
// @Failure 500 {object} walletResponse.GetTraderHistoryErrorResponse
// @failure 502 {object} walletResponse.GetTraderHistoryErrorResponse
// @Router /wallets/{traderID}/history [get]
func (h *WalletHandler) GetTraderHistory(c *gin.Context) {
	traderID := c.Param("traderID")
	slog.Info("trader history", "traderID", traderID)
	if traderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "traderID path param required"})
		return
	}

	proxyResp, err := http.Get("http://localhost:3000/wallets/"+traderID+"/history")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	proxyRespBody, err := io.ReadAll(proxyResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read wallet-service response body"})
		return
	}

	var response walletResponse.GetTraderHistoryResponse
	if err := json.Unmarshal(proxyRespBody, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "json unmarshal failed"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get trader crypto balance
// @Description Get trader crypto balance
// @Tags wallets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param traderID path string true "TraderID"
// @Success 200 {object} walletResponse.GetTraderBalanceResponse
// @Failure 400 {object} walletResponse.GetTraderBalanceErrorResponse
// @Failure 500 {object} walletResponse.GetTraderBalanceErrorResponse
// @Failure 502 {object} walletResponse.GetTraderBalanceErrorResponse
// @Router /wallets/{traderID}/balance [get] 
func (h *WalletHandler) GetTraderBalance(c *gin.Context) {
	traderID := c.Param("traderID")
	slog.Info("trader history", "traderID", traderID)
	if traderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "traderID path param required"})
		return
	}

	proxyResp, err := http.Get("http://localhost:3000/wallets/"+traderID+"/balance")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	proxyRespBody, err := io.ReadAll(proxyResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read wallet-service response body"})
		return
	}

	var response walletResponse.GetTraderBalanceResponse
	if err := json.Unmarshal(proxyRespBody, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "json unmarshal failed"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get trader crypto wallet address
// @Description Get trader crypto wallet address
// @Tags wallets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param traderID path string true "TraderID"
// @Success 200 {object} walletResponse.GetTraderWalletAddressResponse
// @Failure 400 {object} walletResponse.GetTraderWalletAddressErrorResponse
// @Failure 500 {object} walletResponse.GetTraderWalletAddressErrorResponse
// @Failure 502 {object} walletResponse.GetTraderWalletAddressErrorResponse
// @Router /wallets/{traderID}/address [get] 
func (h WalletHandler) GetTraderWalletAddress(c *gin.Context) {
	traderID := c.Param("traderID")
	slog.Info("trader history", "traderID", traderID)
	if traderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "traderID path param required"})
		return
	}

	proxyResp, err := http.Get("http://localhost:3000/wallets/"+traderID+"/address")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	proxyRespBody, err := io.ReadAll(proxyResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read wallet-service response body"})
		return
	}

	var response walletResponse.GetTraderWalletAddressResponse
	if err := json.Unmarshal(proxyRespBody, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "json unmarshal failed"})
		return
	}

	c.JSON(http.StatusOK, response)
}