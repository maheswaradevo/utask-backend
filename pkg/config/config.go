package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	Database      Database
	ServerAddress string
	WhiteListed   string
	ApiSecretKey  string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleStateString  string
}

type Database struct {
	Username string
	Password string
	Address  string
	Port     string
	Name     string
}

var config Config

func Init() {
	err := godotenv.Load()

	if err != nil {
		log.Printf("ERROR .env Not found")
	}

	config.ServerAddress = os.Getenv("SERVER_ADDRESS")
	config.Port = os.Getenv("PORT")
	config.Database.Username = os.Getenv("DB_USERNAME")
	config.Database.Password = os.Getenv("DB_PASSWORD")
	config.Database.Address = os.Getenv("DB_ADDRESS")
	config.Database.Port = os.Getenv("DB_PORT")
	config.Database.Name = os.Getenv("DB_NAME")
	config.WhiteListed = os.Getenv("WHITELISTED_URLS")

	config.GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	config.GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	config.GoogleStateString = os.Getenv("GOOGLE_STATE_STRING")
}

func GetDatabase(username, password, address, databaseName string) *sql.DB {
	log.Printf("INFO GetDatabase database connection: starting database connection process")

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		username, password, address, databaseName)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("Error GetDatabase sql open connection fatal error: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("ERROR GetDatabase db ping fatal error: %v", err)
	}
	log.Printf("INFO GetDatabase database connectionn: established successfully\n")
	return db
}

func GetConfig() *Config {
	return &config
}
