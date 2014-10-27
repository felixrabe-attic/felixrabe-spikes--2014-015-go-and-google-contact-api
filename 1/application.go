package main

import (
	"fmt"
	"os"

	"./google"
)

func main() {
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	if clientId == "" || clientSecret == "" {
		fmt.Println("Usage: CLIENT_ID=xxx CLIENT_SECRET=xxx go run application.go")
		os.Exit(1)
	}

	authCode, err := google.GetAuthCode(clientId)
	if err != nil {
		fatal(err)
	}
	// fmt.Println(authCode)

	accessToken, refreshToken, err := google.RequestTokens(clientId, clientSecret, authCode)
	if err != nil {
		fatal(err)
	}
	fmt.Println(authCode)
	fmt.Println()
	fmt.Println(accessToken)
	fmt.Println()
	fmt.Println(refreshToken)
}

func fatal(x interface{}) {
	fmt.Println(x)
	os.Exit(1)
}
