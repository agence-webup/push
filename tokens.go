package push

import "fmt"

// Platform is an enum representing supported Push platforms
type Platform int

const (
	// IOS represents APNS platform
	IOS Platform = 1
	// Android represents GCM platform
	Android Platform = 2
)

// Token represents a token associated to a uuid (typically a user)
type Token struct {
	UUID     string   `json:"uuid"`
	Value    string   `json:"token"`
	Platform Platform `json:"platform"`
	Language string   `json:"language"`
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
	if t.Language == "" {
		return fmt.Errorf("'language' is required")
	}
	return nil
}

// TokenRepository defines the behavior for interacting with tokens
type TokenRepository interface {
	FindToken(token Token) (*Token, error)
	GetTokensForUUID(uuid string) ([]Token, error)
	RemoveToken(token Token) (*Token, error)
	SaveToken(t Token) error
}

// TokenBag stores a token list
type TokenBag struct {
	Tokens []Token
}

func (b *TokenBag) AddToken(token Token) error {
	if b.Tokens == nil {
		b.Tokens = []Token{}
	}

	b.Tokens = append(b.Tokens, token)
	return nil
}

func (b *TokenBag) GetTokens() []Token {
	if b.Tokens == nil {
		return []Token{}
	}

	return b.Tokens
}

func (b *TokenBag) ResetTokens() error {
	b.Tokens = []Token{}
	return nil
}
