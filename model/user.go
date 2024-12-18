package model

type User struct {
	User_Id       int    `gorm:"primaryKey"`
	Email	 	string `gorm:"unique"`
	Name 	string `gorm:"unique"`
	Password string `gorm:"not null"`
}