package http

import (
	"net/http"
	"webup/push"
	"webup/push/apns/sideshow"

	"github.com/labstack/echo"
)

type SendResource struct {
	Pusher push.Pusher
}

func (s *SendResource) Send() echo.HandlerFunc {
	return func(c echo.Context) error {
		s.Pusher = new(sideshow.Pusher)
		// s.Pusher.Setup()

		tokens := []push.Token{
			push.Token{
				UUID:     "tartanpion",
				Value:    "8ddea0121cfb9c770d5dd593c5332368a82cef5ab8aa71b9ce2a8218ce541166",
				Platform: push.IOS,
				Language: "fr",
			},
		}
		s.Pusher.Send(tokens)

		return c.String(http.StatusOK, "Tokens sent successfully")
	}
}
