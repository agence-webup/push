package http

import (
	"fmt"
	"net/http"
	"webup/push"
	"webup/push/android/fcm"
	"webup/push/apns/sideshow"
	"webup/push/repository/memory"
	"webup/push/service"

	"github.com/labstack/echo"
)

type SendResource struct {
	SendService push.SendService
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

		apnsPusher := new(sideshow.Pusher)
		apnsPusher.Setup()
		fcmPusher := new(fcm.Pusher)
		fcmPusher.Setup()

		s.SendService = service.SendService{
			APNSPusher:      apnsPusher,
			FCMPusher:       fcmPusher,
			TokenRepository: new(memory.TokenRepository),
		}

		s.SendService.Send(notification)

		return c.String(http.StatusOK, "Tokens sent successfully")
	}
}
