package mysql

import (
	"fmt"
	"log"
	"webup/push"
)

type TokenRepository struct {
	Config push.RuntimeConfig
}

func (r *TokenRepository) GetTokens() []push.Token {
	db := GetDB(*r.Config.MySQL)
	tokens := []push.Token{}

	rows, err := db.Query("SELECT uuid, value, platform, language FROM `push_tokens`")
	if err != nil {
		log.Println("Unable to fetch tokens from MySQL", err)
		return tokens
	}
	defer rows.Close()

	for rows.Next() {
		token := push.Token{}
		err := rows.Scan(&token.UUID, &token.Value, &token.Platform, &token.Language)
		if err != nil {
			log.Println("Unable to bind MySQL row to a Token: ", err)
			continue
		}
		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		log.Println("Unable to fetch tokens from MySQL: ", err)
		return tokens
	}

	return tokens
}

func (r *TokenRepository) SetTokens(newTokens []push.Token) {
	log.Println("'SetTokens' not supported in MySQL repository")
}

func (r *TokenRepository) FindToken(token push.Token) (*push.Token, error) {
	db := GetDB(*r.Config.MySQL)

	foundToken := push.Token{}

	err := db.
		QueryRow("SELECT uuid, value, platform, language FROM `push_tokens` WHERE value = ? AND platform = ?", token.Value, token.Platform).
		Scan(&foundToken.UUID, &foundToken.Value, &foundToken.Platform, &foundToken.Language)

	if err != nil {
		return nil, fmt.Errorf("Unable to find token from MySQL: %v", err.Error())
	}

	return &foundToken, nil
}

func (r *TokenRepository) GetTokensForUUID(uuid string) ([]push.Token, error) {
	db := GetDB(*r.Config.MySQL)
	tokens := []push.Token{}

	rows, err := db.Query("SELECT uuid, value, platform, language FROM `push_tokens` where uuid = ?", uuid)
	if err != nil {
		return nil, fmt.Errorf("Unable to fetch tokens by UUID from MySQL: %v", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		token := push.Token{}
		err := rows.Scan(&token.UUID, &token.Value, &token.Platform, &token.Language)
		if err != nil {
			log.Println("Unable to bind MySQL row to a Token: ", err)
			continue
		}
		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		return tokens, fmt.Errorf("Unable to fetch tokens by UUID from MySQL: %v", err.Error())
	}

	return tokens, nil
}

func (r *TokenRepository) RemoveToken(token push.Token) (*push.Token, error) {
	db := GetDB(*r.Config.MySQL)

	foundToken, _ := r.FindToken(token)
	if foundToken == nil {
		return nil, nil
	}

	_, err := db.Exec("DELETE FROM `push_tokens` WHERE value = ? AND platform = ?", token.Value, token.Platform)
	if err != nil {
		return nil, fmt.Errorf("Unable to remove a token from MySQL: %v", err.Error())
	}

	return foundToken, nil
}

func (r *TokenRepository) SaveToken(t push.Token) error {
	db := GetDB(*r.Config.MySQL)

	existingToken, _ := r.FindToken(t)
	if existingToken != nil {
		r.RemoveToken(*existingToken)
	}

	_, err := db.Exec("INSERT INTO `push_tokens` (`uuid`, `value`, `platform`, `language`, `created_at`) VALUES (?,?,?,?, NOW())", t.UUID, t.Value, t.Platform, t.Language)
	if err != nil {
		return fmt.Errorf("Unable to save a token into MySQL: %v", err.Error())
	}

	return nil
}
