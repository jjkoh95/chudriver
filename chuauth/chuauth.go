package chuauth

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2/google"
)

// GetClientFromJSON - this should get config client
func GetClientFromJSON(credentialFile, scope string) *http.Client {
	b, err := ioutil.ReadFile(credentialFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.JWTConfigFromJSON(b, scope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return config.Client(context.Background())
}
