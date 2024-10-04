package events

import "github.com/google/go-github/v65/github"

type PushEventHandler struct {
	Event github.PushEvent
}
