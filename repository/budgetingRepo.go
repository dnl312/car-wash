package repository

import (
	"car-wash/model"

	"gorm.io/gorm"
)

type BudgetingRepo struct {
	DB *gorm.DB
}

func NewBudgetingRepo(db *gorm.DB) *BudgetingRepo {
	return &BudgetingRepo{DB: db}
}

func (r *BudgetingRepo) GetTopupTempByOrderID(order_id string, user_id int) (model.TopUpTemp, error) {
	var topUpTemp model.TopUpTemp
	err := r.DB.Table("topup_temp_p2w4").Where("order_id = ? AND user_id = ?", order_id, user_id).First(&topUpTemp).Error
	if err != nil {
		return topUpTemp, err
	}
	return topUpTemp, nil
}

func (r *BudgetingRepo) TopUpBalance(user_id int, amount string) error {
	err := r.DB.Table("users_l3p2w4").Where("user_id = ?", user_id).Update("balance", gorm.Expr("balance + ?", amount)).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *BudgetingRepo) InsertIntoTopUpTemp(order_id string, user_id int) error {
	topUpTemp := model.TopUpTemp{
		Order_ID: order_id,
		User_ID: user_id,
	}
	err := r.DB.Table("topup_temp_p2w4").Create(&topUpTemp).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *BudgetingRepo) UpdateTopUpTemp(order_id string) error {
	err := r.DB.Table("topup_temp_p2w4").Where("order_id = ?", order_id).Update("status", "SETTLEMENT").Error
	if err != nil {
		return err
	}
	return nil
}

func (r *BudgetingRepo) GetTransactionByUserID(user_id int) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.DB.Table("transactions_p2w4").Where("user_id = ?", user_id).Order("transaction_id desc").Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}	