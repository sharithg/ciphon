package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/sharithg/siphon/internal/remote"
	"golang.org/x/crypto/ssh"
)

func (app *Application) getWorkflows(w http.ResponseWriter, r *http.Request) {
	workflows, err := app.Store.WorkflowRunsStore.GetWorkflowRuns()

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, workflows); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *Application) triggerWorkflow(w http.ResponseWriter, r *http.Request) {
	workflowId := chi.URLParam(r, "workflowId")

	// ssh := app.

	workflows, err := app.Store.WorkflowRunsStore.GetById(workflowId)

	if len(workflows) == 0 {
		app.badRequestResponse(w, r, errors.New("invalid workflow id"))
		return
	}

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	nodes, err := app.Store.Nodes.All()

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	node := nodes[0]

	sshConn, err := remote.New(node.Host, node.User, []byte(node.PemFile), true)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	client, err := sshConn.Dial()

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	gitUrl := workflows[0].Url
	repoName := workflows[0].RepoName
	commitSha := workflows[0].CommitSHA

	for _, workflow := range workflows {
		switch stepType := workflow.Type; stepType {
		case "checkout":
			err := app.Checkout(r.Context(), client, node.PemFile, gitUrl, repoName, commitSha)
			if err != nil {
				app.internalServerError(w, r, err)
				return
			}

			err = app.PullDockerImage(r.Context(), client, workflow.Docker)
			if err != nil {
				app.internalServerError(w, r, err)
				return
			}
		case "restore_cache":
			fmt.Printf("step %s not implemented\n", stepType)
		case "save_cache":
			fmt.Printf("step %s not implemented\n", stepType)
		default:
			cmd := fmt.Sprintf(`
				docker run -v /home/ubuntu/%s:/app/%s %s sh -c "%s"
			`, repoName, repoName, workflow.Docker, workflow.Command)
			err = app.RunStepCommand(r.Context(), client, cmd)
			if err != nil {
				app.internalServerError(w, r, err)
				return
			}
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, "ok"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *Application) RunStepCommand(ctx context.Context, client *ssh.Client, command string) error {
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

func (app *Application) Checkout(ctx context.Context, client *ssh.Client, pemContent, gitUrl, name, ref string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session for checkout: %w", err)
	}
	defer session.Close()
	token, err := remote.GenerateJWTToken([]byte(pemContent), app.Config.Github.AppConfig.OAuth.ClientID)
	if err != nil {
		return err
	}

	cloneUrl, err := convertGitHubURL(gitUrl, token)
	if err != nil {
		return err
	}

	command := fmt.Sprintf(`
        pwd
        if [ -d "%s" ]; then
            rm -rf %s
        fi
        git clone %s
		ls
    `, name, name, cloneUrl)

	if err := remote.RunCommand(session, command); err != nil {
		return err
	}

	return nil
}

func (app *Application) PullDockerImage(ctx context.Context, client *ssh.Client, docker string) error {
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
