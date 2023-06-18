package pkg

import (
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	authHTTPDelivery "github.com/maheswaradevo/utask-backend/internal/authentications/delivery"
	authRedisRepo "github.com/maheswaradevo/utask-backend/internal/authentications/repository/redis"
	authService "github.com/maheswaradevo/utask-backend/internal/authentications/service"
	"github.com/maheswaradevo/utask-backend/internal/notification"
	"github.com/maheswaradevo/utask-backend/pkg/common/helpers"
	"github.com/maheswaradevo/utask-backend/pkg/config"

	calendarHTTPDelivery "github.com/maheswaradevo/utask-backend/internal/calendar/delivery"
	calendarRestRepo "github.com/maheswaradevo/utask-backend/internal/calendar/repository/rest"
	calendarSvc "github.com/maheswaradevo/utask-backend/internal/calendar/service"
)

func Init(router *echo.Echo, rc *redis.Client, cfg config.Config, logger *zap.Logger) {
	app := router.Group("")
	{
		InitAuthModule(app, rc)
		InitCalendarModule(app, cfg, rc, logger)
	}

}

func InitAuthModule(routerGroup *echo.Group, rc *redis.Client) *echo.Group {
	authRedisRepo := authRedisRepo.NewGoogleOauthRedisRepository(rc)
	authService := authService.NewGoogleOauthService(authRedisRepo)
	return authHTTPDelivery.AuthenticationNewDelivery(routerGroup, authService)
}

func InitCalendarModule(routerGroup *echo.Group, cfg config.Config, rc *redis.Client, logger *zap.Logger) *echo.Group {
	restyClient := resty.New().SetDebug(true)
	wavecellClient := notification.NewClient(restyClient.R().
		EnableTrace().
		SetAuthToken(helpers.Env("SMS_API_KEY")).
		SetHeader("Content-Type", "application/json"))
	authRedisRepo := authRedisRepo.NewGoogleOauthRedisRepository(rc)
	calendarRestRepository := calendarRestRepo.NewCalendarRestRepository(logger)
	calendarService := calendarSvc.NewCalendarService(&calendarRestRepository, authRedisRepo, wavecellClient, logger)
	return calendarHTTPDelivery.CalendarNewDelivery(routerGroup, calendarService)
}
