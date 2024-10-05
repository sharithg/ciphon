package workflow

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/sharithg/siphon/internal/remote"
	"github.com/sharithg/siphon/internal/storage"
	"golang.org/x/crypto/ssh"
)

type WorkflowManager struct {
	store          *storage.Storage
	GithubClientId string
}

func New(s *storage.Storage, c string) *WorkflowManager {
	return &WorkflowManager{store: s, GithubClientId: c}
}

func (wm *WorkflowManager) TriggerWorkflow(ctx context.Context, workflowId string) error {

	workflows, err := wm.store.WorkflowRunsStore.GetById(workflowId)
	if err != nil {
		return err
	}

	if len(workflows) == 0 {
		return err
	}

	nodes, err := wm.store.Nodes.All()
	if err != nil {
		return err
	}

	node := nodes[0]
	sshConn, err := remote.New(node.Host, node.User, []byte(node.PemFile), true)
	if err != nil {
		return err
	}

	client, err := sshConn.Dial()
	if err != nil {
		return err
	}

	gitUrl := workflows[0].Url
	repoName := workflows[0].RepoName
	ref := workflows[0].Branch

	defer func() {
		go wm.stopDockerImage(client, repoName)
	}()

	for _, workflow := range workflows {
		switch stepType := workflow.Type; stepType {
		case "checkout":
			if err := wm.pullDockerImage(client, workflow.Docker); err != nil {
				return err
			}

			if err := wm.runBackgroundDockerImage(client, workflow.Docker, repoName); err != nil {
				return err
			}

			if err := wm.checkout(client, node.PemFile, gitUrl, repoName, ref); err != nil {
				return err
			}

		case "restore_cache":
			fmt.Printf("step %s not implemented\n", stepType)
		case "save_cache":
			fmt.Printf("step %s not implemented\n", stepType)
		default:
			cmd := fmt.Sprintf(`
				docker exec -w /%s %s sh -c "%s"
			`, repoName, repoName, workflow.Command)
			if err := wm.runStepCommand(client, cmd); err != nil {
				return err
			}
		}
	}

	return nil
}

func (wm *WorkflowManager) runStepCommand(client *ssh.Client, command string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session for docker pull: %w", err)
	}
	defer session.Close()

	if err := remote.RunCommand(session, command); err != nil {
		return err
	}

	return nil
}

func (wm *WorkflowManager) checkout(client *ssh.Client, pemContent, gitUrl, name, ref string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session for checkout: %w", err)
	}
	defer session.Close()

	token, err := remote.GenerateJWTToken([]byte(pemContent), wm.GithubClientId)
	if err != nil {
		return err
	}

	cloneUrl, err := convertGitHubURL(gitUrl, token)
	if err != nil {
		return err
	}

	command := fmt.Sprintf(`
    docker exec %s sh -c "
        pwd
        git clone %s
        cd %s && git fetch origin && git checkout %s
		git log -1
    "
	`, name, cloneUrl, name, ref)

	if err := remote.RunCommand(session, command); err != nil {
		return err
	}

	return nil
}

func (wm *WorkflowManager) pullDockerImage(client *ssh.Client, docker string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session for docker pull: %w", err)
	}
	defer session.Close()

	command := fmt.Sprintf(`
        docker pull %s
    `, docker)

	if err := remote.RunCommand(session, command); err != nil {
		return err
	}

	return nil
}

func (wm *WorkflowManager) runBackgroundDockerImage(client *ssh.Client, imageName, repoName string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session for docker run: %w", err)
	}
	defer session.Close()

	command := fmt.Sprintf(`
        docker run -d -v /home/ubuntu/%s:/app/%s --name %s %s tail -f /dev/null
    `, repoName, repoName, repoName, imageName)

	if err := remote.RunCommand(session, command); err != nil {
		return err
	}

	return nil
}

func (wm *WorkflowManager) stopDockerImage(client *ssh.Client, repoName string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session for docker stop: %w", err)
	}
	defer session.Close()

	command := fmt.Sprintf(`
		if [ "$(docker ps -q -f name=%s)" ]; then
			docker stop %s
		fi
		if [ "$(docker ps -aq -f name=%s)" ]; then
			docker rm %s
		fi
    `, repoName, repoName, repoName, repoName)

	if err := remote.RunCommand(session, command); err != nil {
		return err
	}

	return nil
}

func convertGitHubURL(originalURL string, token string) (string, error) {
	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse GitHub URL: %v", err)
	}

	pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(pathParts) < 2 {
		return "", fmt.Errorf("invalid GitHub URL format")
	}
	owner, repo := pathParts[0], pathParts[1]

	convertedURL := fmt.Sprintf("https://x-access-token:%s@%s/%s/%s.git", token, parsedURL.Host, owner, repo)
	return convertedURL, nil
}
