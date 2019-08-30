package http

import (
	"fmt"
	"net/http"
	"strconv"
	"webup/push"

	"github.com/labstack/echo"
)

// TokenResource handles tokens requests
type TokenResource struct {
	Repository push.TokenRepository
}

func (r *TokenResource) GetTokens() echo.HandlerFunc {
	return func(c echo.Context) error {
		uuid := c.Param("uuid")
		if uuid == "" {
			return c.String(http.StatusNotFound, "uuid is required")
		}

		tokens, err := r.Repository.GetTokensForUUID(uuid)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("cannot get tokens: %v", err))
		}

		return c.JSON(http.StatusOK, tokens)
	}
}

// AddToken adds a token to a user
// Do nothing if the token is already associated to the user
// Token can be moved between users (a token can only be associated to one user at a time)
func (r *TokenResource) AddToken() echo.HandlerFunc {
	return func(c echo.Context) error {

		token := push.Token{}
		err := c.Bind(&token)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("cannot parse token JSON payload: %v", err))
		}

		if err := token.Validate(); err != nil {
			return c.String(http.StatusUnprocessableEntity, err.Error())
		}

		err = r.Repository.SaveToken(token)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		tokensForUUID, _ := r.Repository.GetTokensForUUID(token.UUID)
		json := map[string]interface{}{
			"tokens": tokensForUUID,
		}

		return c.JSON(http.StatusOK, json)
	}
}

// RemoveToken remove a token
func (r *TokenResource) RemoveToken() echo.HandlerFunc {
	return func(c echo.Context) error {

		rawPlatform, err := strconv.Atoi(c.Param("platform"))
		if err != nil {
			return c.String(http.StatusBadRequest, "platform param is not an int")
		}
		platform := push.Platform(rawPlatform)
		if platform != push.IOS && platform != push.Android {
			return c.String(http.StatusBadRequest, "platform must be 1 (iOS) or 2 (Android)")
		}

		token := push.Token{
			Value:    c.Param("value"),
			Platform: platform,
		}

		removedToken, err := r.Repository.RemoveToken(token)
		if err != nil {
			return c.String(http.StatusInternalServerError, "unable to remove the token")
		}

		if removedToken != nil {
			tokensForUUID, _ := r.Repository.GetTokensForUUID((*removedToken).UUID)
			json := map[string]interface{}{
				"tokens": tokensForUUID,
			}

			return c.JSON(http.StatusOK, json)
		}

		return c.NoContent(http.StatusNotModified)

	}
}
