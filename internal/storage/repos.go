package storage

import (
	"database/sql"
	"time"
)

type CreateRepo struct {
	RepoID        int64
	Name          string
	Owner         string
	Description   string
	URL           string
	RepoCreatedAt time.Time
	RawData       string
}

type ListRepo struct {
	RepoID        int64     `json:"repoId"`
	Name          string    `json:"name"`
	Owner         string    `json:"owner"`
	Description   string    `json:"description"`
	URL           string    `json:"url"`
	RepoCreatedAt time.Time `json:"repoCreatedAt"`
}

type RepoStore struct {
	db *sql.DB
}

func (s *RepoStore) Create(repo CreateRepo) (string, error) {
	var id string

	query := `
	INSERT INTO github_repos (repo_id, name, owner, description, url, repo_created_at, raw_data)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`

	err := s.db.QueryRow(query, repo.RepoID, repo.Name, repo.Owner, repo.Description, repo.URL, repo.RepoCreatedAt, repo.RawData).Scan(&id)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *RepoStore) All() ([]ListRepo, error) {
	var repos []ListRepo

	query := `
	SELECT repo_id, name, owner, description, url, repo_created_at
	FROM github_repos
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var repo ListRepo
		err := rows.Scan(&repo.RepoID, &repo.Name, &repo.Owner, &repo.Description, &repo.URL, &repo.RepoCreatedAt)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return repos, nil
}
