package service

import (
	"car-wash/config"
	internal "car-wash/middleware"
	"car-wash/model"
	"car-wash/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// RentCar godoc
// @Summary Rent a car
// @Description Rent a car by providing car ID and quantity
// @Tags Car
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param transactionRequest body model.TransactionRequest true "Transaction Request"
// @Success 200 {object} map[string]interface{} "success"
// @Failure 400 {object} map[string]string "invalid request"
// @Failure 401 {object} map[string]string "invalid access token"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /car/rent [post]
func RentCar(c echo.Context) error {
	userId, err := internal.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid access token"})
	}

	var transactionRequest model.TransactionRequest
	if err := c.Bind(&transactionRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request", "detail": err.Error()})
	}
	
	err = repository.NewCarRentRepo(config.DB).RentCar(userId, transactionRequest)
	if err != nil {
		if err.Error() == "transaction in progress" {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "transaction in progress"})
		}else if err.Error() == gorm.ErrRecordNotFound.Error() {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "car not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "detail": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{"message": "success", "detail": transactionRequest})
}

// ReturnCar godoc
// @Summary Return a rented car
// @Description Return a rented car by providing transaction ID
// @Tags Car
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param transaction_id path string true "Transaction ID"
// @Success 201 {object} map[string]interface{} "return car success"
// @Failure 400 {object} map[string]string "invalid request"
// @Failure 401 {object} map[string]string "invalid access token"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /car/return/{transaction_id} [post]
func ReturnCar (c echo.Context) error {
	userId, err := internal.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid access token"})
	}
	
	err = repository.NewCarRentRepo(config.DB).ReturnCar(userId,c.Param("transaction_id"))
	if err != nil {
		if err.Error() == "transaction not found" {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "transaction not found"})
		}else if err.Error() == "car not found" {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "car not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "detail": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{"message": "return car success"})
}