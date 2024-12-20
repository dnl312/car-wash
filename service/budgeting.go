package service

import (
	"car-wash/config"
	internal "car-wash/middleware"
	"car-wash/model"
	"car-wash/repository"
	"car-wash/utils"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

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

func TopupSettlement(c echo.Context, cfg *config.MidtransConfig) error {
	userId, err := internal.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid access token"})
	}

	newUrl := cfg.MidtransURLApi+"/v2/" + c.Param("transaction_id")+ "/status"
	headers := map[string]string{
		"authorization":  cfg.MidtransAPIKey,
	}

	res, err := utils.RequestGET(newUrl, headers)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed parsing response from server", "detail": err.Error()})
	}

	var result model.TransactionStatus
	err = json.Unmarshal(res, &result)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "error while unmarshalling response", "detail": err.Error()})
	}

	_, err = repository.NewBudgetingRepo(config.DB).GetTopupTempByOrderID(c.Param("transaction_id"), userId)
	if err != nil {
		if err == gorm.ErrRecordNotFound{
			return c.JSON(http.StatusNotFound, map[string]string{"message": "topup transaction not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "error while getting topup temp"})
	}

	err = repository.NewBudgetingRepo(config.DB).TopUpBalance(userId, result.GrossAmount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "detail": err.Error()})
	}

	err = repository.NewBudgetingRepo(config.DB).UpdateTopUpTemp(c.Param("transaction_id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "detail": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}


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
func CreateTopUp(c echo.Context, cfg *config.MidtransConfig) error {
	userId, err := internal.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid access token"})
	}
	newUrl := cfg.MidtransURL + "/snap/v1/transactions"
	headers := map[string]string{
		"authorization":          cfg.MidtransAPIKey,
		"content-type": "application/json",
	}

	var topUpRequest model.TopUpRequest
	if err := c.Bind(&topUpRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request", "detail": err.Error()})
	}

	currentTime := time.Now()
	timestampString := currentTime.Format("20060102150405")

	requestBody := map[string]interface{}{
		"transaction_details": map[string]interface{}{
			"order_id":    "topup_" + timestampString,
			"gross_amount": topUpRequest.Amount,
		},
		"credit_card": map[string]interface{}{
			"secure": false,
		},
	}
	
	jsonBody, err := json.Marshal(requestBody)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed parsing marshal req body"})
    }

	var midtransResponse model.MidtransResponse
	payload := strings.NewReader(string(jsonBody))
	response, err := utils.RequestPOST(newUrl, headers, payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to make request", "detail": err.Error()})
	}

	if len(response) == 0 {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "empty response from server"})
	}

	err = json.Unmarshal(response, &midtransResponse)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed parsing response from server", "detail": err.Error()})
	}

	repository.NewBudgetingRepo(config.DB).InsertIntoTopUpTemp("topup_" + timestampString, userId)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":        midtransResponse.Token,
		"redirect_url": midtransResponse.RedirectURL,
	})
}