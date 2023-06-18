package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/maheswaradevo/utask-backend/internal/calendar"
	"github.com/maheswaradevo/utask-backend/internal/models"
	"github.com/maheswaradevo/utask-backend/pkg/common"
	"github.com/maheswaradevo/utask-backend/pkg/common/helpers"
)

type CalendarHTTPDelivery struct {
	common.Controller
	routerGroupV1   *echo.Group
	calendarService calendar.CalendarService
}

func CalendarNewDelivery(routerGroupV1 *echo.Group, calendarService calendar.CalendarService) (routeGroup *echo.Group) {
	calendarDelivery := CalendarHTTPDelivery{
		routerGroupV1:   routerGroupV1,
		calendarService: calendarService,
	}
	routeGroup = calendarDelivery.routerGroupV1.Group("calendar")
	{
		routeGroup.GET("/", calendarDelivery.GetEvents)
		routeGroup.GET("/:eventId", calendarDelivery.GetEventByID)
		routeGroup.POST("/messages/:eventId", calendarDelivery.SendSMS)
	}
	return
}

func (c CalendarHTTPDelivery) GetEvents(ctx echo.Context) error {
	eventList, err := c.calendarService.GetEvent(helpers.Context(ctx))
	if err != nil {
		return c.InternalServerError(ctx, &common.APIResponse{
			Code:    500,
			Message: "Internal Server Error",
			Errors:  err,
		})
	}
	return c.Ok(ctx, eventList)
}

func (c CalendarHTTPDelivery) GetEventByID(ctx echo.Context) error {
	eventId := ctx.Param("eventId")
	eventDetail, err := c.calendarService.GetEventByID(helpers.Context(ctx), eventId)
	if err != nil {
		return c.InternalServerError(ctx, &common.APIResponse{
			Code:    500,
			Message: "Internal Server Error",
			Errors:  err,
		})
	}
	return c.Ok(ctx, eventDetail)
}

func (c CalendarHTTPDelivery) SendSMS(ctx echo.Context) error {
	var notifRequest models.NotificationRequest

	if err := ctx.Bind(&notifRequest); err != nil {
		return c.WrapBadRequest(ctx, &common.APIResponse{
			Code:    400,
			Message: "Bad Request",
			Errors:  err,
		})
	}

	eventId := ctx.Param("eventId")

	res, err := c.calendarService.SendSMS(helpers.Context(ctx), eventId, notifRequest.PhoneNumber)
	if err != nil {
		return c.InternalServerError(ctx, &common.APIResponse{
			Code:    500,
			Message: "Internal Server Error",
			Errors:  err,
		})
	}
	return c.Ok(ctx, res)
}
