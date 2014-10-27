package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/felixrabe-go/misc.v0"

	"./google"
)

func main() {
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	if clientId == "" || clientSecret == "" {
		fmt.Println("Usage: CLIENT_ID=xxx CLIENT_SECRET=xxx go run application.go")
		os.Exit(1)
	}

	api, err := google.New(clientId, clientSecret)
	if err != nil {
		misc.Fatal(err)
	}
	// fmt.Printf("%+v\n", api)

	r, err := api.Get("https://www.google.com/m8/feeds/contacts/default/full")
	if err != nil {
		misc.Fatal(err)
	}
	defer r.Body.Close()
	fmt.Printf("%+v\n", r)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		misc.Fatal(err)
	}
	fmt.Printf("\n%s\n", string(b))
}
