package service

import (
	"log"
	"webup/push"
)

type SendService struct {
	TokenRepository push.TokenRepository
	APNSPusher      push.Pusher
}

func (s SendService) Send(notification push.Notification) error {

	for _, uuid := range notification.UUIDs {
		tokens, err := s.TokenRepository.GetTokensForUUID(uuid)
		if err != nil {
			return err
		}

		log.Printf("Tokens: %+v", tokens)

		for _, token := range tokens {
			if token.Platform == push.IOS {
				s.APNSPusher.Send(notification, []push.Token{token})
			} else if token.Platform == push.Android {
				log.Println("Should send Android token ", token.Value)
			}
		}
	}

	return nil
}
