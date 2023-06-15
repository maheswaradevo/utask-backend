package pkg

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"

	authHTTPDelivery "github.com/maheswaradevo/utask-backend/internal/authentications/delivery"
	authRedisRepo "github.com/maheswaradevo/utask-backend/internal/authentications/repository/redis"
	authService "github.com/maheswaradevo/utask-backend/internal/authentications/service"
	"github.com/maheswaradevo/utask-backend/pkg/config"

	calendarHTTPDelivery "github.com/maheswaradevo/utask-backend/internal/calendar/delivery"
	calendarRestRepo "github.com/maheswaradevo/utask-backend/internal/calendar/repository/rest"
	calendarSvc "github.com/maheswaradevo/utask-backend/internal/calendar/service"
)

func Init(router *echo.Echo, rc *redis.Client, cfg config.Config) {
	app := router.Group("")
	{
		InitAuthModule(app, rc)
		InitCalendarModule(app, cfg, rc)
	}

}

func InitAuthModule(routerGroup *echo.Group, rc *redis.Client) *echo.Group {
	authRedisRepo := authRedisRepo.NewGoogleOauthRedisRepository(rc)
	authService := authService.NewGoogleOauthService(authRedisRepo)
	return authHTTPDelivery.AuthenticationNewDelivery(routerGroup, authService)
}

func InitCalendarModule(routerGroup *echo.Group, cfg config.Config, rc *redis.Client) *echo.Group {
	authRedisRepo := authRedisRepo.NewGoogleOauthRedisRepository(rc)

	calendarRestRepository := calendarRestRepo.NewCalendarRestRepository()
	calendarService := calendarSvc.NewCalendarService(&calendarRestRepository, authRedisRepo)
	return calendarHTTPDelivery.CalendarNewDelivery(routerGroup, calendarService)
}
