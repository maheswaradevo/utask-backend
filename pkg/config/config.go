package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	OauthConfGl = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "https://utask-backend-production.up.railway.app/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/calendar.events", "https://www.googleapis.com/auth/calendar"},
		Endpoint:     google.Endpoint,
	}
	OauthStateStringGl = ""
)

type Config struct {
	Port          string
	Database      Database
	ServerAddress string
	WhiteListed   string
	ApiSecretKey  string
	RedisConfig   Redis

	GoogleClientID       string
	GoogleClientSecret   string
	GoogleStateString    string
	GoogleCalendarAPIKey string
}

type Database struct {
	Username string
	Password string
	Address  string
	Port     string
	Name     string
}

type Redis struct {
	URI          string
	Password     string
	Database     int
	MaxActive    int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	MaxConnAge   time.Duration
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
	config.GoogleCalendarAPIKey = os.Getenv("GOOGLE_CALENDAR_API_KEY")
	config.RedisConfig.URI = os.Getenv("REDIS_POOL_URI")
	config.RedisConfig.Password = os.Getenv("REDIS_POOL_PASSWORD")

	intDB, _ := strconv.Atoi(os.Getenv("REDIS_POOL_DB"))
	config.RedisConfig.Database = intDB

	dialValue, _ := time.ParseDuration(os.Getenv("REDIS_POOL_DIAL_TIMEOUT"))
	config.RedisConfig.DialTimeout = dialValue

	readValue, _ := time.ParseDuration(os.Getenv("REDIS_POOL_READ_TIMEOUT"))
	config.RedisConfig.ReadTimeout = readValue

	writeValue, _ := time.ParseDuration(os.Getenv("REDIS_POOL_WRITE_TIMEOUT"))
	config.RedisConfig.WriteTimeout = writeValue

	idleValue, _ := time.ParseDuration(os.Getenv("REDIS_POOL_IDLE_TIMEOUT"))
	config.RedisConfig.IdleTimeout = idleValue

	maxConnAge, _ := time.ParseDuration(os.Getenv("REDIS_POOL_MAX_CONN_AGE"))
	config.RedisConfig.DialTimeout = maxConnAge
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

func NewRedisClient() (*redis.Client, error) {
	var opts = &redis.Options{
		Addr:         config.RedisConfig.URI,
		Password:     config.RedisConfig.Password,
		DB:           config.RedisConfig.Database,
		DialTimeout:  config.RedisConfig.DialTimeout,
		ReadTimeout:  config.RedisConfig.ReadTimeout,
		WriteTimeout: config.RedisConfig.WriteTimeout,
		IdleTimeout:  config.RedisConfig.IdleTimeout,
		MaxConnAge:   config.RedisConfig.MaxConnAge,
	}
	rc := redis.NewClient(opts)
	ctx := context.Background()

	if err := rc.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return rc, nil
}

func NewOauthGoogle() {
	OauthConfGl.ClientID = config.GoogleClientID
	OauthConfGl.ClientSecret = config.GoogleClientSecret
	OauthStateStringGl = config.GoogleStateString
}

func GetConfig() *Config {
	return &config
}
