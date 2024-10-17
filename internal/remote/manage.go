package remote

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sharithg/siphon/internal/config"
)

func generateYamlConfig(token string) (string, error) {
	cfg := config.AgentConfig{
		Token: token,
	}

	val, err := json.Marshal(cfg)

	if err != nil {
		return "", err
	}

	return string(val), nil
}

func (s *SshConn) InstallTools(token string) error {
	start := time.Now()
	client, err := s.Dial()

	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()

	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	config, err := generateYamlConfig(token)

	if err != nil {
		return fmt.Errorf("failed to make agent config: %w", err)
	}

	cmdFile, err := os.ReadFile("./internal/remote/scripts/install_tools.sh")
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	command := strings.Replace(string(cmdFile), "%s", config, -1)

	fmt.Printf("dial took %v\n", time.Since(start))

	if err := RunCommand(session, command, nil); err != nil {
		return err
	}

	log.Printf("Tools installed successfully on %s", s.Host)
	return nil
}
