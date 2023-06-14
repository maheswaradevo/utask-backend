package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/maheswaradevo/utask-backend/internal/authentications/service"
	"github.com/maheswaradevo/utask-backend/pkg/common"
)

type AuthenticationHTTPDelivery struct {
	common.Controller
	routeGroupV1 *echo.Group
}

func AuthenticationNewDelivery(routeGroupV1 *echo.Group) (routeGroup *echo.Group) {
	authenticationDelivery := AuthenticationHTTPDelivery{
		routeGroupV1: routeGroupV1,
	}
	routeGroup = authenticationDelivery.routeGroupV1.Group("/auth")
	{
		routeGroup.GET("/login-gl", authenticationDelivery.HandleGoogleLogin)
		routeGroup.GET("/google/callback", authenticationDelivery.GoogleCallback)
	}
	return
}

func (a AuthenticationHTTPDelivery) HandleGoogleLogin(ctx echo.Context) error {
	service.HandleGoogleLogin(ctx.Response().Writer, ctx.Request())
	return nil
}

func (a AuthenticationHTTPDelivery) GoogleCallback(ctx echo.Context) error {
	service.CallBackFromGoogle(ctx.Response().Writer, ctx.Request())
	return nil
}
