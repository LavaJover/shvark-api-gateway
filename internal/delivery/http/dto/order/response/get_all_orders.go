package response

type GetAllOrdersResponse struct {
	Orders 		[]Order 	`json:"orders"`
	Pagination  Pagination  `json:"pagination"`
}