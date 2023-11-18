package database

import (
	models "cards/internal/models"
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

type Database interface {
	NewDatabaseConnection() (*gorm.DB, error)
}

func NewDatabaseConnection(automigrate bool) (*gorm.DB, error) {
	// godotenv package to get .env variables.
	envFile, err := godotenv.Read("./config/.env")
	if err != nil {
		log.Fatalf("Error loading environment variables file")
	}

	cfg := Config{ 
		Host:     envFile["DATABASEHOST"],
		Port:     envFile["DATABASEPORT"],
		User:     envFile["DATABASEUSER"],
		Password: envFile["DATABASEPASSWORD"],
		DBName:   envFile["DATABASENAME"],
		SSLMode:  envFile["DATABASESSLMODE"],
	}
	// data source name (DSN)
	dsn := cfg.User +":"+ cfg.Password +"@tcp"+ "(" + cfg.Host + ":" + cfg.Port +")/" + cfg.DBName + "?" + "parseTime=true&loc=Local"
	
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Print("err: ",err)
		panic("Database connection failed")
	}
	database.Logger = logger.Default.LogMode(logger.Info)

	if automigrate {
		database.AutoMigrate(&models.Card{})
		if err != nil {
			log.Print("err: ",err)
			panic("Database failed to automigrate users")
		}
		log.Print("Database successfully migrated users")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()		
	
	dbInstance, _ := database.DB()
    err = dbInstance.PingContext(ctx)
    if err != nil {
        return nil, err
    }


	log.Print("Successfully connected to database")

	return database, nil
}