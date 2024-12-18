package main

import (
	"log"
	"net/http"
	"os"

	"car-wash/config"
	"car-wash/service"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
        e := echo.New()

        e.Use(middleware.Logger())
	    e.Use(middleware.Recover())

        err := godotenv.Load()
        if err != nil {
            log.Fatalf("Error loading .env file")
        }

		e.GET("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "test heroku livecode 3")
		})
		e.POST("/users/login", service.LoginUser)
        e.POST("/users/register", service.RegisterUser)

		// a := e.Group("/users/carts")
		// a.Use(internal.CustomJWTMiddleware)
		// a.POST("", service.InsertCart)
		// a.GET("", service.GetCarts)
		// a.DELETE("/:id", service.DeleteCartItem)

		// b := e.Group("/products")
		// b.GET("", service.GetProducts)
		// b.GET("/:id", service.GetProductDetail)


        config.InitDB()
	    defer config.CloseDB()

	    config.ClearPreparedStatements()
    
        port := os.Getenv("PORT")
        if port == "" {
            port = "8080" 
        }

        e.Logger.Fatal(e.Start(":" + port))

    }