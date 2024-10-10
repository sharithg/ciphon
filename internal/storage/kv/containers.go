package storage

import (
	"fmt"

	"github.com/nickalie/fskv"
)

var (
	Prefix = "containers"
)

type ContainersStore struct {
	db *fskv.DB
}

func (s *ContainersStore) Set(id, status string) error {
	key := s.Prefix(id)
	if err := s.db.Set(key, []byte(status)); err != nil {
		return nil
	}
	return nil
}

func (s *ContainersStore) Get(id string) (string, error) {
	key := s.Prefix(id)
	value, err := s.db.Get(key)

	if err != nil {
		return "", err
	}

	return string(value), nil
}

func (s *ContainersStore) Prefix(key string) string {
	return fmt.Sprintf("containers:%s", key)
}
