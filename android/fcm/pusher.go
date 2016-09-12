package fcm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"webup/push"
)

var client *http.Client

const (
	fcmURL    = "https://fcm.googleapis.com/fcm/send"
	serverKey = "AIzaSyCydSplZbeh2LkB_YUc0L91tRKYRU8IbvQ"
)

type gcmRequest struct {
	To           []string         `json:"registration_ids"`
	Priority     string           `json:"priority"`
	Notification *gcmNotification `json:"notification,omitempty"`
	Data         *gcmNotification `json:"data,omitempty"`
}

type gcmNotification struct {
	Body   string                 `json:"body"`
	Title  string                 `json:"title"`
	Custom map[string]interface{} `json:"custom"`
	// Icon  string
	// Color string
}

type Pusher struct {
}

func (p *Pusher) Setup() error {
	if client != nil {
		return nil
	}

	client = &http.Client{}

	return nil
}

func (p *Pusher) Send(notif push.Notification, tokens []push.Token) error {
	if client == nil {
		err := fmt.Errorf("Pusher must be initialized. You must call 'Setup()' before sending notifications")
		log.Fatal(err)
		return err
	}

	gcmReq := gcmRequest{
		Priority: "high",
	}

	if len(notif.Custom) > 0 {
		gcmReq.Data = &gcmNotification{
			Body:   notif.Text,
			Title:  notif.Title,
			Custom: notif.Custom,
		}
	} else {
		gcmReq.Notification = &gcmNotification{
			Body:   notif.Text,
			Title:  notif.Title,
			Custom: notif.Custom,
		}
	}

	for _, token := range tokens {
		gcmReq.To = append(gcmReq.To, token.Value)
	}

	err := p.makeRequest(gcmReq)
	if err != nil {
		log.Println("Error:", err)
		return err
	}

	return nil
}

func (p *Pusher) makeRequest(request gcmRequest) error {
	data, _ := json.Marshal(request)
	fmt.Printf("%+v", string(data))
	req, err := http.NewRequest("POST", fcmURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "key="+serverKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := (*client).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return nil
}
