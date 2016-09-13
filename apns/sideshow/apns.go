package sideshow

import (
	"fmt"
	"log"
	"webup/push"

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
		log.Println("Cert Error:", pemErr)
		return pemErr
	}

	if p.Config.Sandbox {
		client = apns2.NewClient(cert).Development()
	} else {
		client = apns2.NewClient(cert).Production()
	}

	return nil
}

func (p *Pusher) Send(notif push.Notification, tokens []push.Token) error {
	if len(tokens) == 0 {
		return nil
	}

	if client == nil {
		err := fmt.Errorf("Pusher must be initialized. You must call 'Setup()' before sending notifications")
		log.Println(err)
		return err
	}

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
			log.Println("APNs Error:", err)
			return err
		}

		log.Println("APNs ID:", res.ApnsID)
	}

	return nil
}
