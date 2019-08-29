package postgres

import (
	"fmt"
	"log"
	"webup/push"
)

type TokenRepository struct {
	Config push.RuntimeConfig
}

func (r *TokenRepository) GetTokens() []push.Token {
	db := GetDB(*r.Config.Postgres)
	tokens := []push.Token{}

	rows, err := db.Query("SELECT uuid, value, platform, language FROM " + getTableName(*r.Config.Postgres))
	if err != nil {
		log.Println("Unable to fetch tokens from Postgres", err)
		return tokens
	}
	defer rows.Close()

	for rows.Next() {
		token := push.Token{}
		err := rows.Scan(&token.UUID, &token.Value, &token.Platform, &token.Language)
		if err != nil {
			log.Println("Unable to bind Postgres row to a Token: ", err)
			continue
		}
		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		log.Println("Unable to fetch tokens from Postgres: ", err)
		return tokens
	}

	return tokens
}

func (r *TokenRepository) SetTokens(newTokens []push.Token) {
	log.Println("'SetTokens' not supported in Postgres repository")
}

func (r *TokenRepository) FindToken(token push.Token) (*push.Token, error) {
	db := GetDB(*r.Config.Postgres)

	foundToken := push.Token{}

	err := db.
		QueryRow("SELECT uuid, value, platform, language FROM "+getTableName(*r.Config.Postgres)+" WHERE value = $1 AND platform = $2", token.Value, token.Platform).
		Scan(&foundToken.UUID, &foundToken.Value, &foundToken.Platform, &foundToken.Language)

	if err != nil {
		return nil, fmt.Errorf("Unable to find token from Postgres: %v", err.Error())
	}

	return &foundToken, nil
}

func (r *TokenRepository) GetTokensForUUID(uuid string) ([]push.Token, error) {
	db := GetDB(*r.Config.Postgres)
	tokens := []push.Token{}

	rows, err := db.Query("SELECT uuid, value, platform, language FROM "+getTableName(*r.Config.Postgres)+" where uuid = $1", uuid)
	if err != nil {
		return nil, fmt.Errorf("Unable to fetch tokens by UUID from Postgres: %v", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		token := push.Token{}
		err := rows.Scan(&token.UUID, &token.Value, &token.Platform, &token.Language)
		if err != nil {
			log.Println("Unable to bind Postgres row to a Token: ", err)
			continue
		}
		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		return tokens, fmt.Errorf("Unable to fetch tokens by UUID from Postgres: %v", err.Error())
	}

	return tokens, nil
}

func (r *TokenRepository) RemoveToken(token push.Token) (*push.Token, error) {
	db := GetDB(*r.Config.Postgres)

	foundToken, _ := r.FindToken(token)
	if foundToken == nil {
		return nil, nil
	}

	_, err := db.Exec("DELETE FROM "+getTableName(*r.Config.Postgres)+" WHERE value = $1 AND platform = $2", token.Value, token.Platform)
	if err != nil {
		return nil, fmt.Errorf("Unable to remove a token from Postgres: %v", err.Error())
	}

	return foundToken, nil
}

func (r *TokenRepository) SaveToken(t push.Token) error {
	db := GetDB(*r.Config.Postgres)

	existingToken, _ := r.FindToken(t)
	if existingToken != nil {
		r.RemoveToken(*existingToken)
	}

	_, err := db.Exec("INSERT INTO "+getTableName(*r.Config.Postgres)+" (uuid, value, platform, language, created_at) VALUES ($1,$2,$3,$4, NOW())", t.UUID, t.Value, t.Platform, t.Language)
	if err != nil {
		return fmt.Errorf("Unable to save a token into Postgres: %v", err.Error())
	}

	return nil
}
