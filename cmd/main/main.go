package main

import (
	"fmt"

	"github.com/maheswaradevo/utask-backend/internal/authentications/service"
	"github.com/maheswaradevo/utask-backend/pkg"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/maheswaradevo/utask-backend/pkg/config"
)

func main() {
	config.Init()

	cfg := config.GetConfig()

	// db := config.GetDatabase(cfg.Database.Username, cfg.Database.Password, cfg.Database.Address, cfg.Database.Name)
	app := echo.New()

	pkg.Init(app)

	address := fmt.Sprintf("%s:%s", "0.0.0.0", cfg.Port)

	service.InitializeOauthGoogle()

	app.Start(address)
}
