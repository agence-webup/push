package push

import "fmt"

type Pusher interface {
	Setup() error
	Send(notification Notification, tokens []Token) (SendResponse, error)
}

type Notification struct {
	UUIDs  []string               `json:"uuids"`
	Title  string                 `json:"title"`
	Text   string                 `json:"text"`
	Custom map[string]interface{} `json:"custom"`
}

type SendResponse struct {
	InvalidTokens []Token
}

func (n Notification) Validate() error {
	if len(n.UUIDs) == 0 {
		return fmt.Errorf("'uuids' is required and must contains at least 1 UUID")
	}
	if n.Text == "" {
		return fmt.Errorf("'text' is required")
	}

	return nil
}

type SendService interface {
	Send(notification Notification) error
}
