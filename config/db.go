package config

import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "github.com/joho/godotenv"
    "os"
)

var DB *gorm.DB

func GetGormDB() (*gorm.DB, error) {
    if DB != nil {
        return DB, nil
    }

    // Load environment variables from .env file
    err := godotenv.Load(".env")
    if err != nil {
        return nil, err
    }

    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")

    dsn := "user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " host=" + dbHost + " port=" + dbPort + " sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    DB = db

    return DB, nil
}

func GetDB() *gorm.DB {
    return DB
}
