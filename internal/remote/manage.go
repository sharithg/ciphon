package remote

import (
	"fmt"
	"log"
	"time"
)

func (s *SshConn) InstallTools() error {
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

	command := `
		sudo apt-get update && \
			sudo apt-get install -y docker.io git

		if ! getent group docker >/dev/null; then
			sudo groupadd docker
		fi

		if ! groups $USER | grep &>/dev/null '\bdocker\b'; then
			sudo usermod -aG docker $USER
		fi

		newgrp docker
	`
	fmt.Printf("dial took %v\n", time.Since(start))

	if err := RunCommand(session, command); err != nil {
		return err
	}

	log.Printf("Tools installed successfully on %s", s.Host)
	return nil
}
