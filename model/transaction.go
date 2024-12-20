package model

type Transaction struct {
	Transaction_Id int `gorm:"primaryKey"`
	User_Id int `gorm:"not null"`
	Car_Id int `gorm:"not null"`
	Quantity int `gorm:"not null"`
	Total_Amount float64 `gorm:"not null"`
	Start_Date string `gorm:"not null"`
	End_Date string `gorm:"not null"`
	Status string `gorm:"not null"`
}

type TransactionRequest struct {
	Car_Id int `json:"car_id" validate:"required"`
	Quantity int `json:"quantity" validate:"required"`
	Start_Date string `json:"start_date" validate:"required"`
	End_Date string `json:"end_date" validate:"required"`
}

type TransactionStatus struct {
		StatusCode             string `json:"status_code"`
		TransactionID          string `json:"transaction_id"`
		GrossAmount            string `json:"gross_amount"`
		Currency               string `json:"currency"`
		OrderID                string `json:"order_id"`
		PaymentType            string `json:"payment_type"`
		SignatureKey           string `json:"signature_key"`
		TransactionStatus      string `json:"transaction_status"`
		FraudStatus            string `json:"fraud_status"`
		StatusMessage          string `json:"status_message"`
		MerchantID             string `json:"merchant_id"`
		TransactionTime        string `json:"transaction_time"`
		SettlementTime         string `json:"settlement_time"`
		ExpiryTime             string `json:"expiry_time"`
		ChannelResponseCode    string `json:"channel_response_code"`
		ChannelResponseMessage string `json:"channel_response_message"`
		Bank                   string `json:"bank"`
		ApprovalCode           string `json:"approval_code"`
		MaskedCard             string `json:"masked_card"`
		CardType               string `json:"card_type"`
		OnUs                   bool   `json:"on_us"`
	}