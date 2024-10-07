package remote

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

var Green = "\033[32m"
var Reset = "\033[0m"

type SshConn struct {
	Host string
	Key  []byte
	conf *ssh.ClientConfig
}

type streamFunc func(streamType string, buf []byte)

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

func (s *SshConn) Dial() (*ssh.Client, error) {
	conn, err := ssh.Dial("tcp", s.Host, s.conf)

	if err != nil {
		return nil, fmt.Errorf("error connecting to server: %w", err)
	}

	return conn, nil
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

func streamOutput(name string, pipe io.Reader, fn streamFunc) {
	buf := make([]byte, 1024)
	for {
		n, err := pipe.Read(buf)
		if n > 0 {
			if fn != nil {
				fn(name, buf[:n])
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error reading from %s: %v", name, err)
			break
		}
	}
}

func RunCommand(session *ssh.Session, command string, fn streamFunc) error {

	fmt.Println(Green + command + Reset)

	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to set up stdout for command: %w", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to set up stderr for command: %w", err)
	}

	stdoutReader := bufio.NewReader(stdout)
	stderrReader := bufio.NewReader(stderr)

	if err := session.Start(command); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Create channels to read from stdout and stderr
	stdoutChan := make(chan string)
	stderrChan := make(chan string)
	doneChan := make(chan struct{})

	go func() {
		for {
			line, err := stdoutReader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					log.Printf("Error reading from stdout: %v", err)
				}
				close(stdoutChan)
				return
			}
			stdoutChan <- line
		}
	}()

	go func() {
		for {
			line, err := stderrReader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					log.Printf("Error reading from stderr: %v", err)
				}
				close(stderrChan)
				return
			}
			stderrChan <- line
		}
	}()

	// Read from stdout and stderr channels in order
	go func() {
		for stdoutChan != nil || stderrChan != nil {
			select {
			case line, ok := <-stdoutChan:
				if !ok {
					stdoutChan = nil
				} else {
					if fn != nil {
						fn("stdout", []byte(line))
					}
				}
			case line, ok := <-stderrChan:
				if !ok {
					stderrChan = nil
				} else {
					if fn != nil {
						fn("stderr", []byte(line))
					}
				}
			}
		}

		close(doneChan)
	}()

	if err := session.Wait(); err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	<-doneChan

	return nil
}

func GenerateJWTToken(pemContent []byte, clientID string) (string, error) {
	signingKey, err := jwt.ParseRSAPrivateKeyFromPEM(pemContent)
	if err != nil {
		return "", err
	}

	now := time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": now,
		"exp": now + 600,
		"iss": clientID,
	})

	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
