package repository

import (
	"car-wash/model"
	"errors"

	"gorm.io/gorm"
)

type CarRentRepo struct {
	DB *gorm.DB
}

func NewCarRentRepo(db *gorm.DB) *CarRentRepo {
	return &CarRentRepo{DB: db}
}

func (r *CarRentRepo) RentCar(user_id int, transactionReq model.TransactionRequest) error {
	car, err := r.GetCar(transactionReq.Car_Id)
	if err != nil {
		return gorm.ErrRecordNotFound
	}

	var transaction model.Transaction
	err = r.DB.Table("transactions_p2w4").Where("user_id = ? AND car_id = ? AND status = ?", user_id, transactionReq.Car_Id,"IN PROGRESS").First(&transaction).Error
	if err != gorm.ErrRecordNotFound && err != nil {
		return err
	}else if transaction.Transaction_Id != 0 {
		return errors.New("transaction in progress")
	}

	totalAmount := car.Cost * float64(transactionReq.Quantity)

	transaction = model.Transaction{	
		User_Id: user_id, 
		Car_Id: transactionReq.Car_Id, 
		Quantity: transactionReq.Quantity, 
		Total_Amount: totalAmount,
		Start_Date: transactionReq.Start_Date,
		End_Date: transactionReq.End_Date,
		Status: "IN PROGRESS",
	}
	
	err = r.DB.Table("transactions_p2w4").Create(&transaction).Error
	if err != nil {
		return err
	}

	err = r.DB.Table("users_l3p2w4").Where("user_id = ?", user_id).Update("balance", gorm.Expr("balance - ?", totalAmount)).Error
	if err != nil {
		return err
	}

	err = r.DB.Table("cars_p2w4").Where("car_id = ?", transactionReq.Car_Id).Update("quantity", gorm.Expr("quantity - ?", transactionReq.Quantity)).Error
	if err != nil {
		return err
	}
	
	return nil
}

func (r *CarRentRepo) ReturnCar(user_id int, Transaction_Id string) error {
	var transaction model.Transaction
	err := r.DB.Table("transactions_p2w4").Where("user_id = ? AND transaction_id = ? AND status = ? ", user_id, Transaction_Id, "IN PROGRESS").First(&transaction).Error
	if err != nil {
		return err
	}

	_, err = r.GetCar(transaction.Car_Id)
	if err != nil {
		return err
	}

		err = r.DB.Table("cars_p2w4").Where("car_id = ?", transaction.Car_Id).Update("quantity", gorm.Expr("quantity + ?", transaction.Quantity)).Error
	if err != nil {
		return err
	}


	err = r.DB.Table("transactions_p2w4").Where("user_id = ? AND transaction_id = ? ", user_id, Transaction_Id).Update("status", "COMPLETED").Error
	if err != nil {
		return err
	}

	return  nil
}

func (r *CarRentRepo) GetCar(carId int) (*model.Car, error) {
    var car model.Car
    err := r.DB.Table("cars_p2w4").Where("car_id = ? AND quantity > 0", carId).First(&car).Error
    if err != nil {
        return nil, err
    }
    return &car, nil
}