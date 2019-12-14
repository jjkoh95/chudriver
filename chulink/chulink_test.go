package chulink_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jjkoh95/chudriver/chulink"
	"google.golang.org/api/firebasedynamiclinks/v1"
	option "google.golang.org/api/option"
)

const (
	firebaseRegisteredDomain = "chudriver.page.link"
	longURL                  = "https://github.com/jjkoh95/chudriver"
)

var chulinkWrapper chulink.Chulink

func TestMain(m *testing.M) {
	login()
	code := m.Run()
	os.Exit(code)
}

func login() {
	clientOptions := []option.ClientOption{
		option.WithCredentialsFile("credentials.json"),
	}
	chulinkWrapper.Link, _ = firebasedynamiclinks.NewService(context.Background(), clientOptions...)
}

func TestGenerateLink(t *testing.T) {
	_, err := chulinkWrapper.GenerateLink(longURL, firebaseRegisteredDomain)
	if err != nil {
		fmt.Println(err)
		t.Error("Should create url without error")
	}
}
