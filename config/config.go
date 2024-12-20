package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type MidtransConfig struct {
	MidtransAPIKey string
	MidtransURL    string
	MidtransURLApi    string
}

func LoadMidtransConfig() *MidtransConfig {
	return &MidtransConfig{
		MidtransAPIKey: os.Getenv("MIDTRANS_APIKEY"),
		MidtransURL:    os.Getenv("MIDTRANS_URL"),
		MidtransURLApi:    os.Getenv("MIDTRANS_URL_API"),
	}
}

func InitDB() {
    dsn := os.Getenv("DATABASE_URL")
    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
        PrepareStmt: false, 
    })
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }
}

func CloseDB() {
    sqlDB, err := DB.DB()
    if err != nil {
        log.Fatalf("Error getting sql.DB from gorm.DB: %v", err)
    }
    if err := sqlDB.Close(); err != nil {
        log.Fatalf("Error closing database connection: %v", err)
    }
}

func ClearPreparedStatements() {
    sqlDB, err := DB.DB()
    if err != nil {
        log.Fatalf("Error getting sql.DB from gorm.DB: %v", err)
    }
    _, err = sqlDB.Exec("DEALLOCATE ALL")
    if err != nil {
        log.Fatalf("Error deallocating prepared statements: %v", err)
    }
}