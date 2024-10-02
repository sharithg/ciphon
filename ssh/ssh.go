package ssh

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type SshConn struct {
	Host string
	Key  []byte
	conf *ssh.ClientConfig
}

func New(host string, user string, pemKeyContent []byte, ignoreHostKey bool) (*SshConn, error) {
	signer, err := ssh.ParsePrivateKey(pemKeyContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	knownHostsPath, err := knownHostsPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get known_hosts path: %w", err)
	}

	var hostKeyCallback ssh.HostKeyCallback
	if ignoreHostKey {
		hostKeyCallback = ssh.InsecureIgnoreHostKey()
	} else {
		hostKeyCallback, err = knownhosts.New(knownHostsPath)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize known_hosts: %w", err)
		}
	}

	conf := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: hostKeyCallback,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		BannerCallback: func(message string) error {
			fmt.Printf("Banner: %s\n", message)
			return nil
		},
	}

	return &SshConn{
		Host: fmt.Sprintf("%s:22", host),
		Key:  pemKeyContent,
		conf: conf,
	}, nil
}

func (s *SshConn) Ping() error {
	fmt.Println("Connecting to:", s.Host, s.conf.User)
	conn, err := ssh.Dial("tcp", s.Host, s.conf)
	if err != nil {
		return fmt.Errorf("error connecting to server: %w", err)
	}

	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("error creating SSH session: %w", err)
	}
	defer session.Close()

	fmt.Println("SSH connection established successfully.")
	return nil
}

func knownHostsPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get user home directory: %w", err)
	}

	return homeDir + "/.ssh/known_hosts", nil
}
