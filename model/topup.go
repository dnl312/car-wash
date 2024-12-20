package model

type TopUpRequest struct {
	Amount float64 `json:"amount" form:"amount" query:"amount"`
}

type TopUpTemp struct {
    Order_ID string    `gorm:"column:order_id;not null"`
    User_ID  int       `gorm:"column:user_id;not null"`
}

type MidtransResponse struct {
    Token       string `json:"token"`
    RedirectURL string `json:"redirect_url"`
}