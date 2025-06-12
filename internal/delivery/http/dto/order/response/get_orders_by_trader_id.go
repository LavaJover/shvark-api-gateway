package response

type GetOrdersByTraderIDResponse struct {
	Orders []Order `json:"orders"`
}