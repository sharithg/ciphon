package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v65/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
)

type CommitHandler struct {
	githubapp.ClientCreator

	Preamble string
}

func (h *CommitHandler) Handles() []string {
	return []string{"pull_request"}
}

func (h *CommitHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.PullRequest

	fmt.Println("handling event")
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse issue comment event payload")
	}

	fmt.Println(string(payload))

	return nil
}
