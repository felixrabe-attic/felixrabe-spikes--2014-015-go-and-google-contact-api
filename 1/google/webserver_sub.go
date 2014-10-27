// +build ignore

package main

import (
	"fmt"
	// "io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

var result = make(chan string)

func main() {
	switch os.Args[1] {
	case "auth":
		fmt.Println(runAuthWebServer())
		fmt.Println(<-result)
	case "token":
		fmt.Println(runTokenWebServer())
		fmt.Println(<-result)
	}
}

func runAuthWebServer() (port int) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	port = l.Addr().(*net.TCPAddr).Port
	http.HandleFunc("/", authHandler)
	go func(l net.Listener) {
		if err := http.Serve(l, nil); err != nil {
			log.Fatal(err)
		}
	}(l)
	return
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You can close this window now."))
	v := r.URL.Query()
	// ioutil.WriteFile("/tmp/http-request.txt", []byte(fmt.Sprintf("%q\n", r)), 0666)
	if errString := v.Get("error"); errString != "" {
		result <- "error: " + errString
	} else {
		result <- "code: " + v.Get("code")
	}
}

func runTokenWebServer() (port int) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	port = l.Addr().(*net.TCPAddr).Port
	http.HandleFunc("/", tokenHandler)
	go func(l net.Listener) {
		if err := http.Serve(l, nil); err != nil {
			log.Fatal(err)
		}
	}(l)
	return
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You can close this window now."))
	v := r.URL.Query()
	// ioutil.WriteFile("/tmp/http-request.txt", []byte(fmt.Sprintf("%q\n", r)), 0666)
	if errString := v.Get("error"); errString != "" {
		result <- "error: " + errString
	} else {
		result <- "code: " + v.Get("code")
	}
}
