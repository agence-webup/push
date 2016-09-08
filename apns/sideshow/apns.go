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
}

func (p *Pusher) Setup() error {
	cert, pemErr := certificate.FromPemFile("../../cert.pem", "")
	if pemErr != nil {
		log.Println("Cert Error:", pemErr)
		return pemErr
	}

	client = apns2.NewClient(cert).Development()

	return nil
}

func (p *Pusher) Send(tokens []push.Token) error {
	if client == nil {
		err := fmt.Errorf("Pusher must be initialized. You must call 'Setup()' before sending notifications")
		log.Fatal(err)
		return err
	}

	for _, token := range tokens {
		notification := &apns2.Notification{}
		notification.DeviceToken = token.Value
		notification.Topic = "com.ymage.dressinbox"

		payload := payload.NewPayload().AlertBody("youpi")
		notification.Payload = payload

		res, err := client.Push(notification)

		if err != nil {
			log.Println("Error:", err)
			return err
		}

		log.Println("APNs ID:", res.ApnsID)
	}

	return nil
}
