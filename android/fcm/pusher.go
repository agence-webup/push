package fcm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"webup/push"

	log "github.com/Sirupsen/logrus"
)

const (
	maxTokensPerRequest = 1000
)

var client *http.Client

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

type gcmResponse struct {
	Results []gcmResult `json:"results"`
}

type gcmResult struct {
	Error string `json:"error"`
}

type Pusher struct {
	Config push.FCMConfig
}

func NewPusher(config push.FCMConfig) *Pusher {
	pusher := Pusher{
		Config: config,
	}
	return &pusher
}

func (p *Pusher) Setup() error {
	if client != nil {
		return nil
	}

	client = &http.Client{}

	return nil
}

func (p *Pusher) Send(notif push.Notification, tokens []push.Token) (push.SendResponse, error) {
	if len(tokens) == 0 {
		return push.SendResponse{}, nil
	}

	if client == nil {
		err := fmt.Errorf("Pusher must be initialized. You must call 'Setup()' before sending notifications")
		log.Println(err)
		return push.SendResponse{}, err
	}

	// FCM limits the tokens to 1000 per request, so calculate the number of requests to execute
	iterations := int(math.Ceil(float64(len(tokens)) / float64(maxTokensPerRequest)))
	sendResponse := push.SendResponse{}

	for i := 0; i < iterations; i++ {

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

		lowerBound := i * maxTokensPerRequest
		upperBound := (i + 1) * maxTokensPerRequest

		// get the tokens
		selectedTokens := []push.Token{}
		if len(tokens)-lowerBound > maxTokensPerRequest {
			selectedTokens = tokens[lowerBound:upperBound]
		} else {
			selectedTokens = tokens[lowerBound:]
		}

		for _, token := range selectedTokens {
			gcmReq.To = append(gcmReq.To, token.Value)
		}

		// execute the HTTP request
		resp, err := p.makeRequest(gcmReq)
		if err != nil {
			log.WithFields(log.Fields{"error": err, "count": len(selectedTokens)}).Errorln("Unable to execute FCM request")
			continue
		}

		// check if an error is returned for each token
		for i := range selectedTokens {
			result := resp.Results[i]
			if result.Error == "NotRegistered" {
				sendResponse.InvalidTokens = append(sendResponse.InvalidTokens, selectedTokens[i])
			}
		}

	}

	return sendResponse, nil
}

func (p *Pusher) makeRequest(request gcmRequest) (gcmResponse, error) {
	data, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", p.Config.URL, bytes.NewBuffer(data))
	if err != nil {
		return gcmResponse{}, err
	}
	req.Header.Set("Authorization", "key="+p.Config.ServerKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := (*client).Do(req)
	if err != nil {
		return gcmResponse{}, err
	}
	defer resp.Body.Close()

	log.WithFields(log.Fields{"status": resp.StatusCode}).Debugln("FCM request executed")

	body, _ := ioutil.ReadAll(resp.Body)
	log.WithFields(log.Fields{"body": string(body)}).Debugln("FCM request response body")

	gcmResp := gcmResponse{}
	err = json.Unmarshal(body, &gcmResp)
	if err != nil {
		return gcmResponse{}, err
	}

	return gcmResp, nil
}
