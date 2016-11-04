package sideshow

import (
	"fmt"
	"webup/push"

	log "github.com/Sirupsen/logrus"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
)

var client *apns2.Client

type Pusher struct {
	Config push.APNSConfig
}

func NewPusher(config push.APNSConfig) *Pusher {
	pusher := Pusher{
		Config: config,
	}
	return &pusher
}

func (p *Pusher) Setup() error {
	if client != nil {
		return nil
	}

	cert, pemErr := certificate.FromPemFile(p.Config.CertPath, p.Config.CertPass)
	if pemErr != nil {
		log.WithFields(log.Fields{"error": pemErr}).Errorln("APNs certificate error")
		return pemErr
	}

	if p.Config.Sandbox {
		client = apns2.NewClient(cert).Development()
	} else {
		client = apns2.NewClient(cert).Production()
	}

	return nil
}

func (p *Pusher) Send(notif push.Notification, tokens []push.Token) (push.SendResponse, error) {
	if len(tokens) == 0 {
		return push.SendResponse{}, nil
	}

	if client == nil {
		err := fmt.Errorf("Pusher must be initialized. You must call 'Setup()' before sending notifications")
		log.WithFields(log.Fields{"error": err}).Errorln("APNs Pusher not initialized")
		return push.SendResponse{}, err
	}

	invalidTokens := []push.Token{}

	for _, token := range tokens {
		notification := &apns2.Notification{}
		notification.DeviceToken = token.Value
		notification.Topic = p.Config.Topic

		payload := payload.NewPayload().AlertBody(notif.Text)
		if len(notif.Custom) > 0 {
			payload = payload.Custom("custom", notif.Custom)
		}
		notification.Payload = payload

		res, err := client.Push(notification)

		if err != nil {
			if res.Reason == apns2.ReasonDeviceTokenNotForTopic || res.Reason == apns2.ReasonUnregistered {
				invalidTokens = append(invalidTokens, token)
			}

			log.WithFields(log.Fields{"error": err, "id": res.ApnsID}).Errorln("APNs not sent")
			continue
		}

		log.WithFields(log.Fields{"id": res.ApnsID}).Debugln("APNs sent successfully")
	}

	return push.SendResponse{InvalidTokens: invalidTokens}, nil
}
