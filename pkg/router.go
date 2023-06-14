package pkg

import (
	"github.com/labstack/echo/v4"
	authHTTPDelivery "github.com/maheswaradevo/utask-backend/internal/authentications/delivery"
)

func Init(router *echo.Echo) {
	app := router.Group("")
	{
		InitAuthModule(app)
	}

}

func InitAuthModule(routerGroup *echo.Group) *echo.Group {
	return authHTTPDelivery.AuthenticationNewDelivery(routerGroup)
}
