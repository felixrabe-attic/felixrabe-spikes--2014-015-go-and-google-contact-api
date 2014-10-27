package server

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/skratchdot/open-golang/open"
)

type Srv struct {
	l    net.Listener
	Errc chan error
	Reqc chan *Req
}

type Req struct {
	URL   *url.URL
	BodyB []byte
}

func Start() (srv *Srv, err error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	srv = &Srv{
		l:    l,
		Errc: make(chan error, 1),
		Reqc: make(chan *Req, 1),
	}
	smux := http.NewServeMux()
	smux.HandleFunc("/", srv.httpHandler)
	smux.Handle("/favicon.ico", http.NotFoundHandler())
	go func() {
		if err := http.Serve(srv.l, smux); err != nil {
			if netOpErr, ok := err.(*net.OpError); ok && netOpErr.Op == "accept" {
				srv.Errc <- nil
				return
			}
			srv.Errc <- err
			return
		}
		srv.Errc <- nil
		return
	}()
	select {
	case err = <-srv.Errc:
		return nil, err
	case <-time.After(200 * time.Millisecond):
	}
	return
}

func (srv *Srv) OpenAuth(clientId string) error {
	v := url.Values{}
	v.Set("response_type", "code")
	v.Set("client_id", clientId)
	v.Set("redirect_uri", srv.RedirectUri())
	// https://developers.google.com/google-apps/contacts/v3/#authorizing_requests_with_oauth_20
	v.Set("scope", "https://www.googleapis.com/auth/contacts.readonly")

	url := "https://accounts.google.com/o/oauth2/auth?" + v.Encode()
	if err := open.Run(url); err != nil {
		return err
	}
	return nil
}

func (srv *Srv) RedirectUri() string {
	return "http://localhost:" + strconv.Itoa(srv.l.Addr().(*net.TCPAddr).Port)
}

func (srv *Srv) Stop() error {
	srv.l.Close()
	err := <-srv.Errc
	if err != nil {
		return err
	}
	return nil
}

func (srv *Srv) httpHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		srv.Errc <- err
		return
	}
	defer r.Body.Close()

	fmt.Fprint(w, "Thank you. You can now close this window.")
	srv.Reqc <- &Req{URL: r.URL, BodyB: b}
}
