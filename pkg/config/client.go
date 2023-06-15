package config

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

func CalendarServiceClient(token *oauth2.Token) *http.Client {
	return OauthConfGl.Client(context.Background(), token)
}
