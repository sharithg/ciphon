package storage

import "database/sql"

type Storage struct {
	Nodes interface {
		All() ([]Node, error)
		GetById(id string) (*Node, error)
		Create(node Node) (string, error)
		UpdateStatus(nodeID string, status string) error
	}
	Repos interface {
		All() ([]ListRepo, error)
		Create(repo CreateRepo) (string, error)
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Nodes: &NodeStore{db},
		Repos: &RepoStore{db},
	}
}
