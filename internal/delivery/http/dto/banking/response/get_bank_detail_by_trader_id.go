package response

type GetBankDetailsByTraderIDResponse struct {
	BankDetails []BankDetail `json:"bank_details"`
}