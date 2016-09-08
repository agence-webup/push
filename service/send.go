package service

import (
	"log"
	"webup/push"
)

type SendService struct {
	TokenRepository push.TokenRepository
	APNSPusher      push.Pusher
	FCMPusher       push.Pusher
}

func (s SendService) Send(notification push.Notification) error {

	for _, uuid := range notification.UUIDs {
		tokens, err := s.TokenRepository.GetTokensForUUID(uuid)
		if err != nil {
			return err
		}

		log.Printf("Tokens: %+v", tokens)

		fcmTokens := []push.Token{}

		for _, token := range tokens {
			if token.Platform == push.IOS {
				s.APNSPusher.Send(notification, []push.Token{token})
			} else if token.Platform == push.Android {
				fcmTokens = append(fcmTokens, token)
				// log.Println("Should send Android token ", token.Value)
			}
		}

		s.FCMPusher.Send(notification, fcmTokens)
	}

	return nil
}
