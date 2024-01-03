package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/cli/oauth"
	"github.com/google/go-github/v57/github"
	"github.com/wolfeidau/action-workflow-check/internal/ptr"
	"github.com/zalando/go-keyring"
)

const (
	// The "Action Workflow Check" OAuth app
	oauthClientID = "3f1747daa84b81e107b1"

	service = "gh-workflow-version-check:github.com"

	defaultConfigPath = "~/.config/gh-workflow-version-check.json"
)

func Login() error {
	flow := &oauth.Flow{
		Host:     oauth.GitHubHost("https://github.com"),
		ClientID: oauthClientID,
		Scopes:   []string{},
	}

	accessToken, err := flow.DeviceFlow()
	if err != nil {
		return fmt.Errorf("failed to use device flow: %w", err)
	}

	client := github.NewClient(nil).WithAuthToken(accessToken.Token)

	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		return fmt.Errorf("failed to get user information: %w", err)
	}

	err = keyring.Set(service, ptr.ToString(user.Login), accessToken.Token)
	if err != nil {
		return fmt.Errorf("failed to set access token: %w", err)
	}

	err = writeConfig(defaultConfigPath, ptr.ToString(user.Login))
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func writeConfig(configPath string, login string) error {
	config := map[string]string{
		"login": login,
	}

	expandedPath := kong.ExpandPath(configPath)

	data, err := json.Marshal(&config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(expandedPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func readConfig(configPath string) (string, error) {
	expandedPath := kong.ExpandPath(configPath)

	data, err := os.ReadFile(expandedPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("you need to login before running a scan to authenticate with GitHub")
		}

		return "", fmt.Errorf("failed to read file: %w", err)
	}

	var config map[string]string

	err = json.Unmarshal(data, &config)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal config: %w", err)
	}

	login, ok := config["login"]
	if !ok {
		return "", fmt.Errorf("failed to read login from config")
	}

	return login, nil
}
