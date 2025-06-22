package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	walletRequest "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/wallet/request"
	walletResponse "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/wallet/response"
)

type HTTPWalletClient struct {

}

func NewHTTPWalletClient() *HTTPWalletClient {
	return &HTTPWalletClient{}
}

func (c *HTTPWalletClient) CreateWallet(traderID string) (string, error) {
	request := walletRequest.CreateWalletRequest{
		TraderID: traderID,
	}

	requestBody, err := json.Marshal(&request)
	if err != nil {
		return "", err
	}

	response, err := http.Post("http://localhost:3000/wallets/create", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		var createResponse walletResponse.CreateWalletResponse
		if err := json.Unmarshal(responseBody, &createResponse); err != nil {
			return "", err
		}
		return createResponse.Address, nil
	}
	return "", fmt.Errorf("failed to create wallet for trader")
}