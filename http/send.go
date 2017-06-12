package http

import (
	"fmt"
	"net/http"
	"webup/push"
	"webup/push/service"

	"github.com/labstack/echo"
)

type SendResource struct {
	Config          push.RuntimeConfig
	SendService     push.SendService
	TokenRepository push.TokenRepository
}

func (s SendResource) Send() echo.HandlerFunc {
	return func(c echo.Context) error {

		notification := push.Notification{}
		err := c.Bind(&notification)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("cannot parse token JSON payload: %v", err))
		}

		if err := notification.Validate(); err != nil {
			return c.String(http.StatusUnprocessableEntity, err.Error())
		}

		s := service.Get(s.Config, s.TokenRepository)
		s.Send(notification)

		return c.NoContent(http.StatusOK)
	}
}
