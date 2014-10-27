package google

import (
	"bufio"
	"io"
	"net/url"
	"os/exec"
	"strconv"
	"strings"

	"github.com/skratchdot/open-golang/open"
	"gopkg.in/felixrabe-go/misc.v0"
)

func GetAuthCode(clientId string) (authCode string, err error) {
	ws, err := startAWebserver()
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
	v.Set("response_type", "code")
	v.Set("client_id", clientId)
	v.Set("redirect_uri", "http://localhost:"+strPort)
	// https://developers.google.com/google-apps/contacts/v3/#authorizing_requests_with_oauth_20
	v.Set("scope", "https://www.googleapis.com/auth/contacts.readonly")

	url := "https://accounts.google.com/o/oauth2/auth?" + v.Encode()
	open.Run(url)
	authCode, err = ws.result()
	return
}

type wsAType struct {
	stdout io.ReadCloser
	r      *bufio.Reader
	cmd    *exec.Cmd
}

func startAWebserver() (ws *wsAType, err error) {
	goFile := misc.ThisDirJoin("webserver_sub.go")
	cmd := exec.Command("go", "run", goFile, "auth")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	if err = cmd.Start(); err != nil {
		return
	}
	return &wsAType{stdout: stdout, r: bufio.NewReader(stdout), cmd: cmd}, nil
}

func (ws *wsAType) port() (port int, err error) {
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

func (ws *wsAType) result() (authCode string, err error) {
	line, err := ws.r.ReadString('\n')
	if err != nil {
		return
	}
	result := line[:len(line)-1]
	resultSplit := strings.SplitN(result, ": ", 2)
	if resultSplit[0] == "error" {
		err = misc.Errorf("%s", resultSplit[1])
	} else {
		authCode = resultSplit[1]
	}
	return
}

func (ws *wsAType) quit() {
	ws.cmd.Wait()
}
