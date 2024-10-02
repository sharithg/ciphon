package ssh

import (
	"time"
)

func InstallTools(ip, user string, privateKey []byte) error {
	time.Sleep(10 * time.Second)
	// sshConfig := &ssh.ClientConfig{
	// 	User: user,
	// 	Auth: []ssh.AuthMethod{
	// 		ssh.PublicKeysCallback(func() ([]ssh.Signer, error) {
	// 			signer, err := ssh.ParsePrivateKey(privateKey)
	// 			if err != nil {
	// 				return nil, err
	// 			}
	// 			return []ssh.Signer{signer}, nil
	// 		}),
	// 	},
	// 	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	// }

	// client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", ip), sshConfig)
	// if err != nil {
	// 	return fmt.Errorf("failed to connect to VM: %w", err)
	// }
	// defer client.Close()

	// session, err := client.NewSession()
	// if err != nil {
	// 	return fmt.Errorf("failed to create SSH session: %w", err)
	// }
	// defer session.Close()

	// var stdout, stderr bytes.Buffer
	// session.Stdout = &stdout
	// session.Stderr = &stderr

	// command := `
	//     sudo apt-get update && \
	//     sudo apt-get install -y docker.io git
	// `

	// if err := session.Run(command); err != nil {
	// 	return fmt.Errorf("failed to install tools: %s, error: %w", stderr.String(), err)
	// }

	// log.Printf("Tools installed successfully on %s", ip)
	return nil
}
