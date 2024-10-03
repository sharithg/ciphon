package storage

import "database/sql"

type Storage struct {
	Nodes interface {
		All() ([]Node, error)
		Create(node Node) (string, error)
		UpdateStatus(nodeID string, status string) error
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Nodes: &NodeStore{db},
	}
}
