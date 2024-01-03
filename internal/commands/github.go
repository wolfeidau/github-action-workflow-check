package commands

import (
	"fmt"

	"github.com/google/go-github/v57/github"
	"github.com/rs/zerolog/log"
	"github.com/zalando/go-keyring"
)

func BuildClient() (client *github.Client, err error) {
	login, err := readConfig(defaultConfigPath)
	if err != nil {
		log.Warn().Msg("no configuration found, login to avoid rate limiting")
	}

	if login != "" {
		accessToken, err := keyring.Get(service, login)
		if err != nil {
			return nil, fmt.Errorf("failed to get access token: %w", err)
		}

		return github.NewClient(nil).WithAuthToken(accessToken), nil
	}

	return github.NewClient(nil), nil
}
