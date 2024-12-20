package service

import (
	"car-wash/config"
	"car-wash/model"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoginUser godoc
// @Summary Login a user
// @Description Login a user by providing email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body model.LoginRequest true "Login Request"
// @Success 200 {object} map[string]interface{} "login success"
// @Failure 400 {object} map[string]string "invalid request parameters"
// @Failure 401 {object} map[string]string "user not found or invalid password"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /auth/login [post]
func LoginUser(c echo.Context) error {
	var req model.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request parameters"})
	}

	var user model.User
	if err := config.DB.Table("users_l3p2w4").Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound{
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "detail": err.Error()})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid password"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.User_Id,
		"exp" : jwt.TimeFunc().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "detail": err.Error()})
	}

	return c.JSON(http.StatusOK, model.LoginResponse{Token: tokenString})
}

func validateRegisterUser(user model.RegisterUser) error {
    validate := validator.New()
    err := validate.Struct(user)
    if err != nil {
        for _, e := range err.(validator.ValidationErrors) {
            switch e.Field() {
            default:
                return fmt.Errorf("%s is %s", e.Field(), e.Tag())
            }
        }
    }
    return nil
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user by providing user details
// @Tags Auth
// @Accept json
// @Produce json
// @Param registerUser body model.RegisterUser true "Register Request"
// @Success 201 {object} map[string]interface{} "register success"
// @Failure 400 {object} map[string]string "invalid request parameters"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /auth/register [post]
func RegisterUser(c echo.Context) error {
	var user model.RegisterUser
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request parameters"})
	}

	err := validateRegisterUser(user)
	if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "detail": err.Error()})
	}	

	newUser := model.User{
		Email: user.Email,
		Name: user.Name,
		Password: string(hashedPassword),
	}

	if err := config.DB.Table("users_l3p2w4").Create(&newUser).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "derroretail": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "success register",
		"user": map[string]string{
            "Name":  newUser.Name,
            "Email": newUser.Email,
        },
	})
}