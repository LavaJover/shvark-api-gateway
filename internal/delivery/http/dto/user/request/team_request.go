package request

type AddTraderToTeamRequest struct {
	TraderID   string  `json:"trader_id" binding:"required,uuid"`
	Commission float64 `json:"commission" binding:"required,min=0,max=100"`
}

type UpdateRelationshipParamsRequest struct {
	Commission float64 `json:"commission" binding:"required,min=0,max=100"`
}