package pkg

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	authHTTPDelivery "github.com/maheswaradevo/utask-backend/internal/authentications/delivery"
	authRedisRepo "github.com/maheswaradevo/utask-backend/internal/authentications/repository/redis"
	authService "github.com/maheswaradevo/utask-backend/internal/authentications/service"
)

func Init(router *echo.Echo, rc *redis.Client) {
	app := router.Group("")
	{
		InitAuthModule(app, rc)
	}

}

func InitAuthModule(routerGroup *echo.Group, rc *redis.Client) *echo.Group {
	authRedisRepo := authRedisRepo.NewGoogleOauthRedisRepository(rc)
	authService := authService.NewGoogleOauthService(authRedisRepo)
	return authHTTPDelivery.AuthenticationNewDelivery(routerGroup, authService)
}
