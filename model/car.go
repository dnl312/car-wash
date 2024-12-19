package model

type Car struct {
	Car_Id       int    `gorm:"primaryKey"`
	Brand	 	string `gorm:"unique"`
	Cost float64 `gorm:"not null"`
	Quantity int `gorm:"not null"`
}

type CarRequest struct {
	Car_Id       int    `json:"car_id" validate:"required"`
	Quantity int `json:"quantity" validate:"required"`
}
