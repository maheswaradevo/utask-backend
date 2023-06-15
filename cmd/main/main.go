package main

import (
	"fmt"
	"log"

	"github.com/maheswaradevo/utask-backend/pkg"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/maheswaradevo/utask-backend/pkg/config"
)

func main() {
	config.Init()

	cfg := config.GetConfig()

	config.NewOauthGoogle()
	// db := config.GetDatabase(cfg.Database.Username, cfg.Database.Password, cfg.Database.Address, cfg.Database.Name)
	app := echo.New()

	rc, err := config.NewRedisClient()
	if err != nil {
		log.Fatalf("error connection to redis: %v", err)
	}

	pkg.Init(app, rc, *cfg)

	address := fmt.Sprintf("%s:%s", "0.0.0.0", cfg.Port)

	app.Start(address)
}
