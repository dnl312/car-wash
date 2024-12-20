package service

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"car-wash/model"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBudgetingRepo struct {
    mock.Mock
}

func (m *MockBudgetingRepo) GetTransactionByUserID(userID uint) ([]model.Transaction, error) {
    args := m.Called(userID)
    transactions, _ := args.Get(0).([]model.Transaction)
    err := args.Error(1)
    return transactions, err
}

func TestGetTransactionByUserID(t *testing.T) {
    e := echo.New()

    mockRepo := new(MockBudgetingRepo)
    
    mockRepo.On("GetTransactionByUserID", uint(1)).Return([]model.Transaction{}, nil)
    mockRepo.On("GetTransactionByUserID", uint(2)).Return([]model.Transaction(nil), fmt.Errorf("internal server error"))
    mockRepo.On("GetTransactionByUserID", uint(3)).Return([]model.Transaction(nil), fmt.Errorf("invalid access token"))

    handler := func(c echo.Context) error {
        userIDParam := c.Param("user_id")
        userID, err := strconv.ParseUint(userIDParam, 10, 32)
        if err != nil {
            return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
        }

        transactions, err := mockRepo.GetTransactionByUserID(uint(userID))
        if err != nil {
            if err.Error() == "internal server error" {
                return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
            }
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid access token"})
        }
        return c.JSON(http.StatusOK, map[string]interface{}{"status": "success", "data": transactions})
    }

    t.Run("Success", func(t *testing.T) {
        req := httptest.NewRequest(http.MethodGet, "/budgeting/transactions", nil)
        req.Header.Set("Authorization", "Bearer YOUR_ACCESS_TOKEN")
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetParamNames("user_id")
        c.SetParamValues("1")
        if assert.NoError(t, handler(c)) {
            assert.Equal(t, http.StatusOK, rec.Code)
            assert.Contains(t, rec.Body.String(), "success")
        }
    })

    t.Run("InternalServerError", func(t *testing.T) {
        req := httptest.NewRequest(http.MethodGet, "/budgeting/transactions", nil)
        req.Header.Set("Authorization", "Bearer YOUR_ACCESS_TOKEN")
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetParamNames("user_id")
        c.SetParamValues("2")
        if assert.NoError(t, handler(c)) {
            assert.Equal(t, http.StatusInternalServerError, rec.Code)
            assert.Contains(t, rec.Body.String(), "internal server error")
        }
    })

    t.Run("Unauthorized", func(t *testing.T) {
        req := httptest.NewRequest(http.MethodGet, "/budgeting/transactions", nil)
        req.Header.Set("Authorization", "Bearer YOUR_ACCESS_TOKEN")
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetParamNames("user_id")
        c.SetParamValues("3")
        if assert.NoError(t, handler(c)) {
            assert.Equal(t, http.StatusUnauthorized, rec.Code)
            assert.Contains(t, rec.Body.String(), "invalid access token")
        }
    })
}