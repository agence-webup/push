package memory

import (
	"fmt"
	"webup/push"
)

type TokenRepository struct {
	Tokens []push.Token
}

func (r *TokenRepository) FindToken(token push.Token) (*push.Token, error) {
	for _, t := range r.Tokens {
		if t.Value == token.Value && t.Platform == token.Platform {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("token not found")
}

func (r *TokenRepository) GetTokensForUUID(uuid string) ([]push.Token, error) {
	tokens := []push.Token{}
	for i, t := range r.Tokens {
		if t.UUID == uuid {
			tokens = append(tokens, r.Tokens[i])
		}
	}

	return tokens, nil
}

func (r *TokenRepository) RemoveToken(token push.Token) (*push.Token, error) {
	for i, t := range r.Tokens {
		if t.Value == token.Value && t.Platform == token.Platform {
			removedToken := t
			r.Tokens = append(r.Tokens[:i], r.Tokens[i+1:]...)
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
	r.Tokens = append(r.Tokens, t)

	return nil
}
