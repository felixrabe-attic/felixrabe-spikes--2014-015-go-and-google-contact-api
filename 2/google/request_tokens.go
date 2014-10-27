package google

import (
	"./token"
)

func (api *API) requestTokens() (tokens *token.Tokens, err error) {
	srv, err := token.StartServer()
	if err != nil {
		return nil, err
	}
	token.SendRequest(api.clientId, api.clientSecret, api.authCode, srv.Port, srv)
	tokens, err = srv.WaitAndClose()

	return
}
