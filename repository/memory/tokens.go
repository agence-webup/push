package memory

import (
	"fmt"
	"webup/push"
)

var tokens *[]push.Token

type TokenRepository struct {
}

func (r *TokenRepository) GetTokens() []push.Token {
	if tokens == nil {
		tokens = new([]push.Token)
	}

	return *tokens
}

func (r *TokenRepository) SetTokens(newTokens []push.Token) {
	tokens = &newTokens
}

func (r *TokenRepository) FindToken(token push.Token) (*push.Token, error) {
	for _, t := range r.GetTokens() {
		if t.Value == token.Value && t.Platform == token.Platform {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("token not found")
}

func (r *TokenRepository) GetTokensForUUID(uuid string) ([]push.Token, error) {
	tokens := []push.Token{}
	for i, t := range r.GetTokens() {
		if t.UUID == uuid {
			tokens = append(tokens, r.GetTokens()[i])
		}
	}

	return tokens, nil
}

func (r *TokenRepository) RemoveToken(token push.Token) (*push.Token, error) {
	for i, t := range r.GetTokens() {
		if t.Value == token.Value && t.Platform == token.Platform {
			removedToken := t
			tokens := append(r.GetTokens()[:i], r.GetTokens()[i+1:]...)
			r.SetTokens(tokens)
			return &removedToken, nil
		}
	}

	return nil, nil
}

func (r *TokenRepository) SaveToken(t push.Token) error {

	existingToken, _ := r.FindToken(t)
	if existingToken != nil {
		r.RemoveToken(*existingToken)
	}
	tokens := append(r.GetTokens(), t)
	r.SetTokens(tokens)

	return nil
}
