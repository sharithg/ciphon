package api

import (
	"fmt"
	"os"
	"strconv"

	"github.com/palantir/go-githubapp/githubapp"
)

func ReadGithubConfig() (*GithubConfig, error) {
	// address := os.Getenv("SERVER_ADDRESS")
	// portStr := os.Getenv("SERVER_PORT")
	integrationIDStr := os.Getenv("GITHUB_INTEGRATION_ID")
	webhookSecret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	privateKey := os.Getenv("GITHUB_PRIVATE_KEY")
	pullRequestPreamble := os.Getenv("APP_PULL_REQUEST_PREAMBLE")
	clientId := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	installationIdStr := os.Getenv("GITHUB_INSTALLATION_ID")

	// port, err := strconv.Atoi(portStr)
	// if err != nil {
	// 	log.Fatalln("error reading gh port", err)
	// }

	integrationID, err := strconv.ParseInt(integrationIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error reading gh integrationID: %w", err)
	}

	installationId, err := strconv.ParseInt(installationIdStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error reading gh installationId: %w", err)
	}

	return &GithubConfig{
		AppConfig: githubapp.Config{
			App: struct {
				IntegrationID int64  "yaml:\"integration_id\" json:\"integrationId\""
				WebhookSecret string "yaml:\"webhook_secret\" json:\"webhookSecret\""
				PrivateKey    string "yaml:\"private_key\" json:\"privateKey\""
			}{
				IntegrationID: integrationID,
				WebhookSecret: webhookSecret,
				PrivateKey:    privateKey,
			},
			OAuth: GithubOAuth{
				ClientID:     clientId,
				ClientSecret: clientSecret,
			},
		},
		InstallationId:      installationId,
		PullRequestPreamble: pullRequestPreamble,
	}, nil
}
