package request

type GetTrafficRequest struct {
	Page int32	`form:"page" binding:"min=1"`
	Limit int32	`form:"limit" binding:"min=1,max=100"`
}