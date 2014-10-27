package token

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"../debug"
)

type answer struct {
	Error        string `json:"error"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
	TokenType    string `json:"token_type"`

	// access_token	The token that can be sent to a Google API.
	// refresh_token	A token that may be used to obtain a new access token, included by default for installed applications. Refresh tokens are valid until the user revokes access.
	// expires_in	The remaining lifetime of the access token.
	// token_type	Identifies the type of token returned. Currently, this field always has the value Bearer.
}

func SendRequest(clientId, clientSecret, authCode string, port int, srv *Srv /* hack */) (err error) {
	// FIRST VARIANT:

	v := url.Values{}
	v.Set("code", authCode)
	v.Set("client_id", clientId)
	v.Set("client_secret", clientSecret)
	v.Set("redirect_uri", "http://localhost:"+strconv.Itoa(port))
	v.Set("grant_type", "authorization_code")
	vString := v.Encode()

	// SECOND VARIANT:

	// vString := "code=" + authCode + "&client_id=" + clientId + "&client_secret=" + clientSecret + "&redirect_uri=http://localhost:" + strconv.Itoa(port) + "&grant_type=authorization_code"

	// THIRD VARIANT:

	// v := url.Values{}
	// v.Set("code", authCode)
	// v.Set("client_id", clientId)
	// v.Set("client_secret", clientSecret)
	// v.Set("redirect_uri", "http://localhost")
	// v.Set("grant_type", "authorization_code")
	// vString := v.Encode()

	// FOURTH VARIANT:

	// vString := "code=" + authCode + "&client_id=" + clientId + "&client_secret=" + clientSecret + "&redirect_uri=urn:ietf:wg:oauth:2.0:oob&grant_type=authorization_code"

	// FIFTH VARIANT:

	// v := url.Values{}
	// v.Set("code", authCode)
	// v.Set("client_id", clientId)
	// v.Set("client_secret", clientSecret)
	// v.Set("redirect_uri", "urn:ietf:wg:oauth:2.0:oob")
	// v.Set("grant_type", "authorization_code")
	// vString := v.Encode()

	debug.Printf("Request: %q", vString)
	body := strings.NewReader(vString)

	url := "https://accounts.google.com/o/oauth2/token"
	resp, err := http.Post(url, "application/x-www-form-urlencoded", body)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	answer := new(answer)
	debug.Printf("Response: %q", string(b))
	err = json.Unmarshal(b, answer)
	if err != nil {
		return
	}
	go func() {
		time.Sleep(3 * time.Second)
		srv.shortcutHack(answer.Error)
	}()
	return
}
