package service

import (
	"car-wash/config"
	internal "car-wash/middleware"
	"car-wash/model"
	"car-wash/repository"
	"net/http"

	"github.com/labstack/echo/v4"
)

// TopUpBalance godoc
// @Summary Top up user balance
// @Description Top up the balance of a user by providing the amount
// @Tags Budgeting
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param amount body float64 true "Amount to top up"
// @Success 201 {object} map[string]interface{} "top up balance success"
// @Failure 400 {object} map[string]string "invalid request"
// @Failure 401 {object} map[string]string "invalid access token"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /budgeting/topup [post]
func TopUpBalance(c echo.Context) error {
	userId, err := internal.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid access token"})
	}

	var topUpRequest model.TopUpRequest
	if err := c.Bind(&topUpRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request", "detail": err.Error()})
	}
	
	err = repository.NewBudgetingRepo(config.DB).TopUpBalance(userId, topUpRequest.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "detail": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{"message": "top up balance success"})
}

// GetTransactionByUserID godoc
// @Summary Get transactions by user ID
// @Description Get all transactions for a user by providing the user ID
// @Tags Budgeting
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]interface{} "success"
// @Failure 401 {object} map[string]string "invalid access token"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /budgeting/transactions [get]
func GetTransactionByUserID (c echo.Context) error {
	userId, err := internal.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid access token"})
	}
	
	transactions, err := repository.NewBudgetingRepo(config.DB).GetTransactionByUserID(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "detail": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "success", "data": transactions})
}
