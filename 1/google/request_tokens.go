package google

import (
	"bufio"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"

	"gopkg.in/felixrabe-go/misc.v0"
)

func RequestTokens(clientId, clientSecret, authCode string) (accessToken, refreshToken string, err error) {
	ws, err := startTWebserver()
	if err != nil {
		return
	}
	defer ws.quit()

	port, err := ws.port()
	if err != nil {
		return
	}
	strPort := strconv.Itoa(port)

	v := url.Values{}
	v.Set("code", authCode)
	v.Set("client_id", clientId)
	v.Set("client_secret", clientSecret)
	v.Set("redirect_uri", "http://localhost:"+strPort)
	v.Set("grant_type", "authorization_code")

	url := "https://accounts.google.com/o/oauth2/token"
	bodyType := "application/x-www-form-urlencoded"
	body := strings.NewReader(v.Encode())
	resp, err := http.Post(url, bodyType, body)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	accessToken = string(b)
	return
}

type wsTType struct {
	stdout io.ReadCloser
	r      *bufio.Reader
	cmd    *exec.Cmd
}

func startTWebserver() (ws *wsTType, err error) {
	goFile := misc.ThisDirJoin("webserver_sub.go")
	cmd := exec.Command("go", "run", goFile, "auth")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	if err = cmd.Start(); err != nil {
		return
	}
	return &wsTType{stdout: stdout, r: bufio.NewReader(stdout), cmd: cmd}, nil
}

func (ws *wsTType) port() (port int, err error) {
	line, err := ws.r.ReadString('\n')
	if err != nil {
		return
	}
	i, err := strconv.Atoi(line[:len(line)-1])
	if err != nil {
		return
	}
	port = i
	return
}

func (ws *wsTType) quit() {
	ws.cmd.Wait()
}
