package auth

import (
	"net/url"
	"strconv"

	"github.com/skratchdot/open-golang/open"
)

func SendRequest(clientId string, port int) {
	v := url.Values{}
	v.Set("response_type", "code")
	v.Set("client_id", clientId)
	v.Set("redirect_uri", "http://localhost:"+strconv.Itoa(port))
	// https://developers.google.com/google-apps/contacts/v3/#authorizing_requests_with_oauth_20
	v.Set("scope", "https://www.googleapis.com/auth/contacts.readonly")

	url := "https://accounts.google.com/o/oauth2/auth?" + v.Encode()
	open.Run(url)
}
