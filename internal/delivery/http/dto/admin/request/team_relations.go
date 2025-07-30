package request

type CreateTeamRelationRequest struct {
	TraderID 			string `json:"traderId"`
	TeamLeadID 			string `json:"teamLeadId"`
	TeamRelationParams 	TeamRelationParams `json:"teamRelationParams"`
}

type TeamRelationParams struct {
	Commission float64 `json:"commission"`
}

type UpdateTeamRelationRequest struct {
	RelationID string `json:"relationId"`
	TeamRelationParams TeamRelationParams `json:"teamRelationParams"`
}