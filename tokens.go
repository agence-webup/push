package push

import "fmt"

type Platform int

const (
	IOS     Platform = 1
	Android Platform = 2
)

// Token represents a token associated to a uuid (typically a user)
type Token struct {
	UUID     string   `json:"uuid"`
	Value    string   `json:"token"`
	Platform Platform `json:"platform"`
}

// Validate returns an error if the Token is not valid
func (t *Token) Validate() error {
	if t.UUID == "" {
		return fmt.Errorf("'uuid' is required")
	}
	if t.Value == "" {
		return fmt.Errorf("'token' is required")
	}
	if t.Platform != IOS && t.Platform != Android {
		return fmt.Errorf("'platform' is required and must be 1 (iOS) or 2 (Android)")
	}
	return nil
}

type TokenRepository interface {
	FindTokenWithValue(value string) (Token, error)
	RemoveTokenWithValue(value string) error
	SaveToken(t Token) error
}
