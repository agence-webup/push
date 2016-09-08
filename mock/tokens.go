package mock

import (
	"fmt"
	"webup/push"
)

type TokenRepository struct {
	Tokens []push.Token
}

func (r *TokenRepository) FindTokenWithValue(value string) (*push.Token, error) {
	for _, t := range r.Tokens {
		if t.Value == value {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("unable to find token")
}

func (r *TokenRepository) RemoveTokenWithValue(value string) error {
	for i := range r.Tokens {
		if t.Value == value {
			r.Tokens = append(r.Tokens[:i], r.Tokens[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("unable to find token")
}

func (r *TokenRepository) SaveToken(t push.Token) error {

	token, _ := r.FindTokenWithValue(t.Value)
	if token != nil {
		r.RemoveTokenWithValue(t.Value)
	}
	r.Tokens = append(r.Tokens, t)

	return nil
}
