package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/LavaJover/shvark-api-gateway/internal/client"
	walletRequest "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/wallet/request"
	walletResponse "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/wallet/response"
	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	WalletClient *client.HTTPWalletClient
}

func NewWalletHandler(walletClient *client.HTTPWalletClient) (*WalletHandler, error) {
	return &WalletHandler{
		WalletClient: walletClient,
	}, nil
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

	proxyResp, err := http.Post(fmt.Sprintf("http://%s/wallets/create", h.WalletClient.Addr), "application/json", bytes.NewBuffer(proxyRequestBody))
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

	if proxyResp.StatusCode >= 200 && proxyResp.StatusCode < 300 {
		var response walletResponse.CreateWalletResponse
		if err := json.Unmarshal(proxyRespBody, &response); err != nil {
			c.Data(proxyResp.StatusCode, "application/json", proxyRespBody)
			return
		}

		c.JSON(proxyResp.StatusCode, response)
		return
	}

	c.Data(proxyResp.StatusCode, "application/json", proxyRespBody)
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

	proxyResp, err := http.Post(fmt.Sprintf("http://%s/wallets/freeze", h.WalletClient.Addr), "application/json", bytes.NewBuffer(proxyRequestBody))
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

	if proxyResp.StatusCode >= 200 && proxyResp.StatusCode < 300 {
		var response walletResponse.FreezeResponse
		if err := json.Unmarshal(proxyRespBody, &response); err != nil {
			c.Data(proxyResp.StatusCode, "application/json", proxyRespBody)
			return
		}

		c.JSON(proxyResp.StatusCode, response)
		return
	}

	c.Data(proxyResp.StatusCode, "application/json", proxyRespBody)
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

	proxyResp, err := http.Post(fmt.Sprintf("http://%s/wallets/release", h.WalletClient.Addr), "application/json", bytes.NewBuffer(proxyRequestBody))
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

	if proxyResp.StatusCode >= 200 && proxyResp.StatusCode < 300 {
		var response walletResponse.ReleaseResponse
		if err := json.Unmarshal(proxyRespBody, &response); err != nil {
			c.Data(proxyResp.StatusCode, "application/json", proxyRespBody)
			return
		}

		c.JSON(proxyResp.StatusCode, response)
		return
	}

	c.Data(proxyResp.StatusCode, "application/json", proxyRespBody)
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

	proxyResp, err := http.Post(fmt.Sprintf("http://%s/wallets/withdraw", h.WalletClient.Addr), "application/json", bytes.NewBuffer(proxyRequestBody))
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


	if proxyResp.StatusCode >= 200 && proxyResp.StatusCode < 300 {
		var response walletResponse.WithdrawResponse
		if err := json.Unmarshal(proxyRespBody, &response); err != nil {
			// Даже если тело успешное, но не парсится — вернём raw JSON
			c.Data(proxyResp.StatusCode, "application/json", proxyRespBody)
			return
		}

		c.JSON(proxyResp.StatusCode, response)
		return
	}

	c.Data(proxyResp.StatusCode, "application/json", proxyRespBody)
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

	proxyResp, err := http.Post(fmt.Sprintf("http://%s/wallets/deposit", h.WalletClient.Addr), "application/json", bytes.NewBuffer(proxyRequestBody))
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

	if proxyResp.StatusCode >= 200 && proxyResp.StatusCode < 300 {
		var response walletResponse.DepositResponse
		if err := json.Unmarshal(proxyRespBody, &response); err != nil {
			c.Data(proxyResp.StatusCode, "application/json", proxyRespBody)
			return
		}

		c.JSON(proxyResp.StatusCode, response)
		return
	}

	c.Data(proxyResp.StatusCode, "application/json", proxyRespBody)
}

// @Summary Withdraw crypto off-chain
// @Description Withdraw crypto off-chain
// @Tags wallets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body walletRequest.OffchainWithdrawRequest true "wallet data"
// @Success 200 {object} walletResponse.OffchainWithdrawResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /wallets/offchain-withdraw [post]
func (h *WalletHandler) OffchainWithdraw(c *gin.Context) {
	var request walletRequest.OffchainWithdrawRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	proxyRequestBody, err := json.Marshal(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "json marshal failed"})
		return
	}

	proxyResp, err := http.Post(fmt.Sprintf("http://%s/wallets/offchain-withdraw", h.WalletClient.Addr), "application/json", bytes.NewBuffer(proxyRequestBody))
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

	if proxyResp.StatusCode >= 200 && proxyResp.StatusCode < 300 {
		var response walletResponse.OffchainWithdrawResponse
		if err := json.Unmarshal(proxyRespBody, &response); err != nil {
			c.Data(proxyResp.StatusCode, "application/json", proxyRespBody)
			return
		}

		c.JSON(proxyResp.StatusCode, response)
		return
	}

	c.Data(proxyResp.StatusCode, "application/json", proxyRespBody)
}

// @Summary Get trader transactions history
// @Description Get trader transaction history with pagination
// @Tags wallets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param traderID path string true "TraderID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
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

	// Формируем URL с передачей параметров пагинации
	walletURL := fmt.Sprintf("http://%s/wallets/%s/history", h.WalletClient.Addr, traderID)
	
	// Переносим query-параметры из оригинального запроса
	queryParams := c.Request.URL.Query()
	if len(queryParams) > 0 {
		walletURL += "?" + queryParams.Encode()
	}

	proxyResp, err := http.Get(walletURL)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer proxyResp.Body.Close() // Важно закрывать тело

	proxyRespBody, err := io.ReadAll(proxyResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read wallet-service response body"})
		return
	}

	if proxyResp.StatusCode >= 200 && proxyResp.StatusCode < 300 {
		var response walletResponse.GetTraderHistoryResponse
		if err := json.Unmarshal(proxyRespBody, &response); err != nil {
			// Логируем ошибку парсинга
			slog.Error("failed to parse wallet-service response", "error", err)
			c.Data(proxyResp.StatusCode, "application/json", proxyRespBody)
			return
		}

		c.JSON(proxyResp.StatusCode, response)
		return
	}

	// Проксируем ошибки как есть
	c.Data(proxyResp.StatusCode, "application/json", proxyRespBody)
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

	proxyResp, err := http.Get(fmt.Sprintf("http://%s/wallets/%s/balance", h.WalletClient.Addr, traderID))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	proxyRespBody, err := io.ReadAll(proxyResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read wallet-service response body"})
		return
	}

	if proxyResp.StatusCode >= 200 && proxyResp.StatusCode < 300 {
		var response walletResponse.GetTraderBalanceResponse
		if err := json.Unmarshal(proxyRespBody, &response); err != nil {
			c.Data(proxyResp.StatusCode, "application/json", proxyRespBody)
			return
		}

		c.JSON(proxyResp.StatusCode, response)
		return
	}

	c.Data(proxyResp.StatusCode, "application/json", proxyRespBody)
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

	proxyResp, err := http.Get(fmt.Sprintf("http://%s/wallets/%s/address", h.WalletClient.Addr, traderID))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	proxyRespBody, err := io.ReadAll(proxyResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read wallet-service response body"})
		return
	}

	if proxyResp.StatusCode >= 200 && proxyResp.StatusCode < 300 {
		var response walletResponse.GetTraderWalletAddressResponse
		if err := json.Unmarshal(proxyRespBody, &response); err != nil {
			c.Data(proxyResp.StatusCode, "application/json", proxyRespBody)
			return
		}

		c.JSON(proxyResp.StatusCode, response)
		return
	}

	c.Data(proxyResp.StatusCode, "application/json", proxyRespBody)
}

// @Summary Get commission profit
// @Description Get total commission profit for user in specified period
// @Tags wallets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param traderID path string true "Trader ID"
// @Param from query string true "Start date (ISO 8601 format)"
// @Param to query string true "End date (ISO 8601 format)"
// @Success 200 {object} walletResponse.CommissionProfitResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /wallets/{traderID}/commission-profit [get]
func (h *WalletHandler) GetCommissionProfit(c *gin.Context) {
    traderID := c.Param("traderID")
    from := c.Query("from")
    to := c.Query("to")

    slog.Info("commission profit request", 
        "traderID", traderID, 
        "from", from, 
        "to", to)

    // Validate parameters
    if traderID == "" || from == "" || to == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "traderID, from and to parameters are required",
        })
        return
    }

    // Create request body for wallet-service
    requestBody := map[string]string{
        "traderId": traderID,
        "from":     from,
        "to":       to,
    }
    jsonBody, _ := json.Marshal(requestBody)

    // Make request to wallet-service
    walletURL := fmt.Sprintf("http://%s/wallets/commission-profit", h.WalletClient.Addr)
    resp, err := http.Post(walletURL, "application/json", bytes.NewBuffer(jsonBody))
    if err != nil {
        slog.Error("wallet-service request failed", "error", err)
        c.JSON(http.StatusBadGateway, gin.H{
            "error": "failed to connect to wallet-service",
        })
        return
    }
    defer resp.Body.Close()

    // Handle wallet-service response
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        slog.Error("failed to read wallet-service response", "error", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "failed to read wallet-service response",
        })
        return
    }

    // Forward response with original status code
    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        var successResp walletResponse.CommissionProfitResponse
        if err := json.Unmarshal(body, &successResp); err == nil {
            c.JSON(resp.StatusCode, successResp)
            return
        }
    }

    // Forward error responses as-is
    c.Data(resp.StatusCode, "application/json", body)
}