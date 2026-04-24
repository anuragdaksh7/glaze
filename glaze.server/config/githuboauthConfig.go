package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var GithubOauthConfig *oauth2.Config

func init() {
	config, err := LoadConfig(".")
	if err != nil {
		panic(err)
	}

	GithubOauthConfig = &oauth2.Config{
		ClientID:     config.GithubClientID,
		ClientSecret: config.GithubClientSecret,
		Endpoint:     github.Endpoint,
	}
}
