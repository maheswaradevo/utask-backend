package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/maheswaradevo/utask-backend/internal/calendar"
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
