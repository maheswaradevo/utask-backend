package main

import (
	"fmt"
	"log"

	"github.com/maheswaradevo/utask-backend/pkg"
	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/maheswaradevo/utask-backend/pkg/config"
)

func main() {
	config.Init()

	cfg := config.GetConfig()

	logger, _ := zap.NewProduction()

	config.NewOauthGoogle()
	// db := config.GetDatabase(cfg.Database.Username, cfg.Database.Password, cfg.Database.Address, cfg.Database.Name)
	app := echo.New()
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	rc, err := config.NewRedisClient()
	if err != nil {
		log.Fatalf("error connection to redis: %v", err)
	}

	pkg.Init(app, rc, *cfg, logger)

	PORT := fmt.Sprintf("%s:%s", "0.0.0.0", cfg.Port)

	app.Start(PORT)
}
