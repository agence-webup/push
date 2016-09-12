package http

import (
	"fmt"
	"net/http"
	"webup/push"
	"webup/push/android/fcm"
	"webup/push/apns/sideshow"
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

		apnsPusher := sideshow.NewPusher(*s.Config.APNS)
		apnsPusher.Setup()

		fcmPusher := fcm.NewPusher(*s.Config.FCM)
		fcmPusher.Setup()

		s.SendService = service.SendService{
			APNSPusher:      apnsPusher,
			FCMPusher:       fcmPusher,
			TokenRepository: s.TokenRepository,
		}

		s.SendService.Send(notification)

		return c.String(http.StatusOK, "Tokens sent successfully")
	}
}
