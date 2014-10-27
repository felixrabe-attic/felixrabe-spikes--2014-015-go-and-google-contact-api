package token

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
)

type Srv struct {
	Port   int
	l      net.Listener
	result chan *result

	closed      bool
	closedMutex sync.Mutex
}

type result struct {
	Error error
	Code  string
}

func StartServer() (srv *Srv, err error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	srv = &Srv{
		Port:   l.Addr().(*net.TCPAddr).Port,
		l:      l,
		result: make(chan *result),
	}
	smux := http.NewServeMux()
	smux.HandleFunc("/", srv.httpHandler)
	smux.Handle("/favicon.ico", http.NotFoundHandler())
	go func() {
		if err := http.Serve(srv.l, smux); err != nil {
			srv.result <- &result{Error: err}
			fmt.Println(err)
		}
	}()
	return
}

func (srv *Srv) WaitAndClose() (tokens *Tokens, err error) {
	res := <-srv.result
	defer srv.l.Close()

	if res.Error != nil {
		return nil, res.Error
	} else {
		return &Tokens{Access: res.Code}, nil
	}
}

func (srv *Srv) shortcutHack(code string) {
	srv.closedMutex.Lock()
	defer srv.closedMutex.Unlock()
	if srv.closed {
		return
	}
	srv.closed = true
	go func() { srv.result <- &result{Code: code} }()
}

func (srv *Srv) httpHandler(w http.ResponseWriter, r *http.Request) {
	srv.closedMutex.Lock()
	defer srv.closedMutex.Unlock()
	if srv.closed {
		return
	}
	srv.closed = true

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		srv.result <- &result{Error: err}
		return
	}
	if len(b) > 0 {
		srv.result <- &result{Code: "BODY: " + string(b)}
		return
	}

	v := r.URL.Query()
	if errString := v.Get("error"); errString != "" {
		srv.result <- &result{Error: errors.New(errString)}
	} else {
		srv.result <- &result{Code: v.Get("code")}
	}
}
