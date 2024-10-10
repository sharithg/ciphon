package config

import (
	"encoding/json"
	"os"
)

type AgentConfig struct {
	Token string `json:"token"`
}

func LoadAgentConfig(path string) (*AgentConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config AgentConfig
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
