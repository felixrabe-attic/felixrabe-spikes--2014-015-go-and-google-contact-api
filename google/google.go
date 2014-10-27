package google

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"./server"
)

type API struct {
	clientId     string
	clientSecret string
	authCode     string
	tokens       *Tokens
	client       *http.Client
}

type Tokens struct {
	Error        string `json:"error"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    uint32 `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func New(clientId, clientSecret string) (api *API, err error) {
	api = &API{
		clientId:     clientId,
		clientSecret: clientSecret,
	}

	srv, err := server.Start()
	if err != nil {
		return
	}
	defer srv.Stop()

	api.authCode, err = api.getAuthCode(srv)
	if err != nil {
		return
	}

	api.tokens, err = api.getTokens(srv)
	if err != nil {
		return
	}

	api.client = &http.Client{}
	return
}

func (api *API) Get(url string) (r *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header["Authorization"] = []string{"Bearer " + api.tokens.AccessToken}
	r, err = api.client.Do(req)
	if err != nil {
		return
	}
	return
}

func (api *API) getAuthCode(srv *server.Srv) (authCode string, err error) {
	err = srv.OpenAuth(api.clientId)
	if err != nil {
		return
	}

	// Wait for incoming request
	var req *server.Req
	select {
	case req = <-srv.Reqc:
	case err = <-srv.Errc:
		return
	}
	v := req.URL.Query()
	if e := v.Get("error"); e != "" {
		err = errors.New(e)
		return
	}
	authCode = v.Get("code")
	return
}

func (api *API) getTokens(srv *server.Srv) (tokens *Tokens, err error) {
	v := url.Values{}
	v.Set("code", api.authCode)
	v.Set("client_id", api.clientId)
	v.Set("client_secret", api.clientSecret)
	v.Set("redirect_uri", srv.RedirectUri())
	v.Set("grant_type", "authorization_code")
	vString := v.Encode()

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

	tokens = new(Tokens)
	err = json.Unmarshal(b, tokens)
	return
}
