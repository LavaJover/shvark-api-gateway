package response

type GetOrdersByTraderIDResponse struct {
	Orders 		[]Order 	`json:"orders"`
	Pagination  Pagination  `json:"pagination"`
}

type Pagination struct {
	CurrentPage  int64 `json:"current_page"`
	TotalPages 	 int64 `json:"total_pages"`
	TotalItems 	 int64 `json:"total_items"`
	ItemsPerPage int64 `json:"items_per_page"`
}