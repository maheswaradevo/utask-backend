package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/maheswaradevo/utask-backend/internal/authentications"
	"github.com/maheswaradevo/utask-backend/pkg/common"
	"github.com/maheswaradevo/utask-backend/pkg/common/helpers"
)

type AuthenticationHTTPDelivery struct {
	common.Controller
	routeGroupV1 *echo.Group
	oauthService authentications.GoogleOauthService
}

func AuthenticationNewDelivery(routeGroupV1 *echo.Group, oauthService authentications.GoogleOauthService) (routeGroup *echo.Group) {
	authenticationDelivery := AuthenticationHTTPDelivery{
		routeGroupV1: routeGroupV1,
		oauthService: oauthService,
	}
	routeGroup = authenticationDelivery.routeGroupV1.Group("/auth")
	{
		routeGroup.GET("/login-gl", authenticationDelivery.HandleGoogleLogin)
		routeGroup.GET("/google/callback", authenticationDelivery.GoogleCallback)
		routeGroup.GET("/logout", authenticationDelivery.Logout)
	}
	return
}

func (a AuthenticationHTTPDelivery) HandleGoogleLogin(ctx echo.Context) error {
	a.oauthService.HandleGoogleLogin(ctx.Response().Writer, ctx.Request())
	return nil
}

func (a AuthenticationHTTPDelivery) GoogleCallback(ctx echo.Context) error {
	t, err := a.oauthService.CallBackFromGoogle(ctx.Response().Writer, ctx.Request())
	if err != nil {
		return a.InternalServerError(ctx, &common.APIResponse{
			Code:    500,
			Message: "Internal Server Error",
			Errors:  err,
		})
	}
	return a.Ok(ctx, t)
}

func (a AuthenticationHTTPDelivery) Logout(ctx echo.Context) error {
	isLogOut, err := a.oauthService.Logout(helpers.Context(ctx))
	if err != nil {
		return a.InternalServerError(ctx, &common.APIResponse{
			Code:    500,
			Message: "Internal Server Error",
			Errors:  err,
		})
	}
	return a.Ok(ctx, isLogOut)
}
