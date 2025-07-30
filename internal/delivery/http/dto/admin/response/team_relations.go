package response

type TeamRelationsResponse struct {
	TeamRelations []TeamRelation `json:"teamRelations"`
}

type TeamRelation struct {
	ID string `json:"id"`
	TraderID string `json:"traderId"`
	TeamLeadID string `json:"teamLeadId"`
	TeamRelationRarams TeamRelationRarams `json:"teamRelationRarams"`
}

type TeamRelationRarams struct {
	Commission float64 `json:"commission"`
}