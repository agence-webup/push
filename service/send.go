package service

import (
	"webup/push"

	log "github.com/Sirupsen/logrus"
)

type SendService struct {
	TokenRepository push.TokenRepository
	APNSPusher      push.Pusher
	FCMPusher       push.Pusher
}

func (s SendService) Send(notification push.Notification) error {

	log.Infoln("Send request received: preparing tokens...")

	apnsTokens := []push.Token{}
	fcmTokens := []push.Token{}

	for _, uuid := range notification.UUIDs {
		tokens, err := s.TokenRepository.GetTokensForUUID(uuid)
		if err != nil {
			return err
		}

		for _, token := range tokens {
			if token.Platform == push.IOS {
				apnsTokens = append(apnsTokens, token)
			} else if token.Platform == push.Android {
				fcmTokens = append(fcmTokens, token)
			}
		}
	}

	log.WithFields(log.Fields{"count": len(apnsTokens)}).Infoln("APNs tokens to send.")
	log.WithFields(log.Fields{"count": len(fcmTokens)}).Infoln("FCM tokens to send.")

	go func() {
		// APNs
		apnsResponse, err := s.APNSPusher.Send(notification, apnsTokens)
		if err != nil {
			log.Println(err)
		}

		// FCM
		fcmResponse, err := s.FCMPusher.Send(notification, fcmTokens)
		if err != nil {
			log.Println(err)
		}

		// clean invalid tokens
		log.WithFields(log.Fields{"count": len(apnsResponse.InvalidTokens)}).Infoln("APNs invalid tokens found.")
		for _, token := range apnsResponse.InvalidTokens {
			s.TokenRepository.RemoveToken(token)
		}
		log.WithFields(log.Fields{"count": len(fcmResponse.InvalidTokens)}).Infoln("FCM invalid tokens found.")
		for _, token := range fcmResponse.InvalidTokens {
			s.TokenRepository.RemoveToken(token)
		}
	}()

	return nil
}
