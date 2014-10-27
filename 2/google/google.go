package google

import (
	"./token"
)

type API struct {
	clientId     string
	clientSecret string
	authCode     string
	tokens       *token.Tokens
}

func New(clientId, clientSecret string) (api *API, err error) {
	api = &API{
		clientId:     clientId,
		clientSecret: clientSecret,
	}

	authCode, err := api.authorize()
	if err != nil {
		return
	}
	api.authCode = authCode

	tokens, err := api.requestTokens()
	if err != nil {
		return
	}
	api.tokens = tokens

	return
}
