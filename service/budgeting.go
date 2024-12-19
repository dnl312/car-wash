package service

import (
	"car-wash/config"
	internal "car-wash/middleware"
	"car-wash/model"
	"car-wash/repository"
	"net/http"

	"github.com/labstack/echo/v4"
)

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
