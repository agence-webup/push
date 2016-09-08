package http

import (
	"fmt"
	"net/http"
	"webup/push"

	"github.com/labstack/echo"
)

// TokenResource handles tokens requests
type TokenResource struct {
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

		return c.String(http.StatusOK, fmt.Sprintf("%+v", token))
	}
}

// RemoveToken remove a token
func (r *TokenResource) RemoveToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Param("token")
		return c.String(http.StatusOK, "try to remove token: "+token)
	}
}
