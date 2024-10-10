package storage

import (
	"log"
	"os"
	"path"

	"github.com/nickalie/fskv"
)

type KvStorage struct {
	Containers interface {
		Set(id, status string) error
		Get(id string) (string, error)
	}
}

func NewKvStorage() *KvStorage {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	path := path.Join(dirname, ".ciphon", "data")

	db, err := fskv.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	return &KvStorage{
		Containers: &ContainersStore{db},
	}
}
