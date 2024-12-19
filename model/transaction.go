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