package github

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/palantir/go-githubapp/githubapp"
)

type HTTPConfig struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type MyApplicationConfig struct {
	PullRequestPreamble string `yaml:"pull_request_preamble"`
}

type Config struct {
	Server HTTPConfig       `yaml:"server"`
	Github githubapp.Config `yaml:"github"`

	AppConfig MyApplicationConfig `yaml:"app_configuration"`

	InstallationId int64
}

func ReadConfig() Config {
	address := os.Getenv("SERVER_ADDRESS")
	portStr := os.Getenv("SERVER_PORT")
	integrationIDStr := os.Getenv("GITHUB_INTEGRATION_ID")
	webhookSecret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	privateKey := os.Getenv("GITHUB_PRIVATE_KEY")
	pullRequestPreamble := os.Getenv("APP_PULL_REQUEST_PREAMBLE")
	clientId := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	installationIdStr := os.Getenv("GITHUB_INSTALLATION_ID")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalln("error reading gh port", err)
	}

	integrationID, err := strconv.ParseInt(integrationIDStr, 10, 64)
	if err != nil {
		log.Fatalln("error reading gh integrationID", err)
	}

	installationId, err := strconv.ParseInt(installationIdStr, 10, 64)
	if err != nil {
		log.Fatalln("error reading gh installationId", err)
	}

	fmt.Println(integrationID)

	return Config{
		Server: HTTPConfig{
			Address: address,
			Port:    port,
		},
		Github: githubapp.Config{
			App: struct {
				IntegrationID int64  `yaml:"integration_id" json:"integrationId"`
				WebhookSecret string `yaml:"webhook_secret" json:"webhookSecret"`
				PrivateKey    string `yaml:"private_key" json:"privateKey"`
			}{
				IntegrationID: integrationID,
				WebhookSecret: webhookSecret,
				PrivateKey:    privateKey,
			},
			OAuth: struct {
				ClientID     string "yaml:\"client_id\" json:\"clientId\""
				ClientSecret string "yaml:\"client_secret\" json:\"clientSecret\""
			}{
				ClientID:     clientId,
				ClientSecret: clientSecret,
			},
			V3APIURL: "https://api.github.com/",
		},
		AppConfig: MyApplicationConfig{
			PullRequestPreamble: pullRequestPreamble,
		},
		InstallationId: installationId,
	}
}
