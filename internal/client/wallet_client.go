package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	walletRequest "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/wallet/request"
	walletResponse "github.com/LavaJover/shvark-api-gateway/internal/delivery/http/dto/wallet/response"
)

type HTTPWalletClient struct {
	Addr string
}

func NewHTTPWalletClient(addr string) *HTTPWalletClient {
	return &HTTPWalletClient{
		Addr: addr,
	}
}

func (c *HTTPWalletClient) CreateWallet(traderID string) (string, error) {
	request := walletRequest.CreateWalletRequest{
		TraderID: traderID,
	}

	requestBody, err := json.Marshal(&request)
	if err != nil {
		return "", err
	}

	response, err := http.Post(fmt.Sprintf("http://%s/wallets/create", c.Addr), "application/json", bytes.NewBuffer(requestBody))
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

type balanceResponse struct {
	UserID  string  `json:"userId"`
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
	Frozen  float64 `json:"frozen"`
}

func (c *HTTPWalletClient) GetBalance(userID string) (float64, error){
	url := fmt.Sprintf("http://%s/wallets/%s/balance",c.Addr, userID)
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result balanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Balance, nil
}


type withdrawRequest struct {
	TraderID  string  `json:"traderId"`
	ToAddress string  `json:"toAddress"`
	Amount    float64 `json:"amount"`
}

type withdrawResponse struct {
	TxHash string `json:"txHash"`
}

func (c *HTTPWalletClient) Withdraw(userID, toAddress string, amount float64) (string, error) {
	reqBody := withdrawRequest{
		TraderID:  userID,
		ToAddress: toAddress,
		Amount:    amount,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(fmt.Sprintf("http://%s/wallets/withdraw", c.Addr), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to make POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var res withdrawResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return res.TxHash, nil
}

type SetWithdrawalRulesRequest struct {
	TraderID 		string 	`json:"traderId"`
	FixedFee 		float64 `json:"fixedFee"`
	MinAmount 		float64 `json:"minAmount"`
	CooldownSeconds int64 	`json:"cooldownSeconds"`
}

type SetWithdrawalRulesResponse struct {
	Success bool 	`json:"success"`
	Rule Rule 		`json:"rule"`
}

type Rule struct {
	ID 				int64 	  `json:"id"`
	TraderID 		string 	  `json:"traderId"`
	FixedFee 		float64   `json:"fixedFee"`
	MinAmount 		float64   `json:"minAmount"`
	CooldownSeconds int64 	  `json:"cooldownSeconds"`
	UpdatedAt 		time.Time `json:"updatedAt"`
	CreatedAt 		time.Time `json:"createdAt"`
}

func (c *HTTPWalletClient) SetWithdrawalRules(requestBody *SetWithdrawalRulesRequest) (*SetWithdrawalRulesResponse, error) {
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(fmt.Sprintf("http://%s/admin/withdrawal-rules", c.Addr), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var res SetWithdrawalRulesResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

type GetWithdrawalRulesResponse struct {
	ID 				int64 	  `json:"id"`
	TraderID 		string 	  `json:"traderId"`
	FixedFee 		float64   `json:"fixedFee"`
	MinAmount 		float64   `json:"minAmount"`
	CooldownSeconds int64 	  `json:"cooldownSeconds"`
	UpdatedAt 		time.Time `json:"updatedAt"`
	CreatedAt 		time.Time `json:"createdAt"`
}

func (c *HTTPWalletClient) GetWithdrawalRules(userID string) (*GetWithdrawalRulesResponse, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/withdrawal-rules/%s", c.Addr, userID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var res GetWithdrawalRulesResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c HTTPWalletClient) DeleteWithdrawalRule(userID string) error {
	url := fmt.Sprintf("http://%s/withdrawal-rules/%s", c.Addr, userID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error deleting withdrawal rule")
	}

	return nil
}