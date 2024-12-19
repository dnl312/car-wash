package model

type TopUpRequest struct {
	Amount float64 `json:"amount" form:"amount" query:"amount"`
}