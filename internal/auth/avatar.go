package auth

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sharithg/siphon/internal/storage/minio"
)

type AvatarStore interface {
	fmt.Stringer
	Put(userID string, reader io.Reader) (avatarID string, err error) // save avatar data from the reader and return base name
	Get(avatarID string) (reader io.ReadCloser, size int, err error)  // load avatar via reader
	ID(avatarID string) (id string)                                   // unique id of stored avatar's data
	Remove(avatarID string) error                                     // remove avatar data
	List() (ids []string, err error)                                  // list all avatar ids
	Close() error                                                     // close store
}

type FileStore struct {
	basePath string
	store    *minio.Storage
}

var _ AvatarStore = &FileStore{}

func NewFileStore(basePath string, store *minio.Storage) *FileStore {
	return &FileStore{basePath: basePath, store: store}
}

func (fs *FileStore) String() string {
	return fmt.Sprintf("FileStore(basePath=%s)", fs.basePath)
}

func (fs *FileStore) Put(userID string, reader io.Reader) (avatarID string, err error) {
	avatarID = userID + ".avatar"
	filePath := filepath.Join(fs.basePath, avatarID)

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return "", err
	}

	return avatarID, nil
}

func (fs *FileStore) Get(avatarID string) (reader io.ReadCloser, size int, err error) {
	ctx := context.Background()
	file, err := os.CreateTemp(fs.basePath, avatarID)
	if err != nil {
		return nil, 0, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, 0, err
	}

	_, err = fs.store.Upload(ctx, "avatars", avatarID, file.Name(), "png")
	if err != nil {
		return nil, 0, err
	}

	return file, int(stat.Size()), nil
}

func (fs *FileStore) ID(avatarID string) string {
	return strings.TrimSuffix(avatarID, ".avatar")
}

func (fs *FileStore) Remove(avatarID string) error {
	ctx := context.Background()
	err := fs.store.DeleteObject(ctx, "avatars", avatarID)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FileStore) List() ([]string, error) {
	ctx := context.Background()
	ids, err := fs.store.ListObjects(ctx, "avatars")
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (fs *FileStore) Close() error {
	return nil
}
