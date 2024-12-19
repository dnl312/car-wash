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

func (r *BudgetingRepo) TopUpBalance(user_id int, amount float64) error {
	err := r.DB.Table("users_l3p2w4").Where("user_id = ?", user_id).Update("balance", gorm.Expr("balance + ?", amount)).Error
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