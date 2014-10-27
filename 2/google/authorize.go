package google

import (
	"./auth"
)

func (api *API) authorize() (authCode string, err error) {
	srv, err := auth.StartServer()
	if err != nil {
		return "", err
	}
	auth.SendRequest(api.clientId, srv.Port)
	authCode, err = srv.WaitAndClose()

	return
}
