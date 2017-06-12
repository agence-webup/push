package service

import (
	"webup/push"
	"webup/push/android/fcm"
	"webup/push/apns/sideshow"

	log "github.com/Sirupsen/logrus"
)

var service *SendService

type SendService struct {
	TokenRepository push.TokenRepository
	Pushers         push.PushManagerByPlatform
}

func Get(config push.RuntimeConfig, tokenRepo push.TokenRepository) SendService {
	if service == nil {

		managers := make(push.PushManagerByPlatform)

		if config.APNS != nil {
			pusher := sideshow.NewPushManager(*config.APNS)
			pusher.Setup()
			managers[push.IOS] = pusher
		} else {
			log.Warnln("APNs config is not present. Cannot send APNs notifications")
		}

		if config.FCM != nil {
			pusher := fcm.NewPushManager(*config.FCM)
			pusher.Setup()
			managers[push.Android] = pusher
		} else {
			log.Warnln("FCM config is not present. Cannot send FCM notifications")
		}

		service = &SendService{
			Pushers:         managers,
			TokenRepository: tokenRepo,
		}
	}

	return *service
}

func (s *SendService) Send(notification push.Notification) error {

	log.Infoln("Send request received: preparing tokens...")

	for _, uuid := range notification.UUIDs {
		tokens, err := s.TokenRepository.GetTokensForUUID(uuid)
		if err != nil {
			return err
		}

		for _, token := range tokens {
			if manager, ok := s.Pushers[token.Platform]; ok {
				manager.AddToken(token)
			}
		}
	}

	for platform, manager := range s.Pushers {
		tokens := manager.GetTokens()
		log.WithFields(log.Fields{"count": len(tokens), "platform": platform}).Infoln("Tokens collected")
	}

	go func() {

		for platform, manager := range s.Pushers {
			response, err := manager.Send(notification)
			if err != nil {
				log.Println(err)
			}

			log.WithFields(log.Fields{"count": len(response.InvalidTokens), "platform": platform}).Infoln("Invalid tokens found")
			for _, token := range response.InvalidTokens {
				s.TokenRepository.RemoveToken(token)
			}
		}

	}()

	return nil
}
