package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sharithg/siphon/internal/auth"
)

type User struct {
	ID         string `json:"id" db:"id"`
	Username   string `json:"username" db:"username"`
	Email      string `json:"email" db:"email"`
	ExternalId string `json:"external_id" db:"external_id"`
	AuthType   string `json:"auth_type" db:"auth_type"`
}

type GitHubUserInfo struct {
	UserID string          `json:"user_id" db:"user_id"`
	Data   auth.GitHubUser `json:"data" db:"data"`
}

type UserDisplay struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatarUrl"`
}

type UserStore struct {
	pool *pgxpool.Pool
}

func (s *UserStore) Create(ctx context.Context, user User, githubInfo GitHubUserInfo) (*User, error) {
	var id string

	userQuery := `
	INSERT INTO users (username, email, external_id, auth_type)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`

	err := s.pool.QueryRow(ctx, userQuery, user.Username, user.Email, user.ExternalId, user.AuthType).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %s", err)
	}

	user.ID = id

	githubQuery := `
	INSERT INTO github_user_info (user_id, data)
	VALUES ($1, $2)
	`

	_, err = s.pool.Exec(ctx, githubQuery, id, githubInfo.Data)
	if err != nil {
		return nil, fmt.Errorf("error creating github_user_info: %s", err)
	}

	return &user, nil
}

func (s *UserStore) GetByExternalId(ctx context.Context, id string) (*User, error) {
	var user User

	query := `
	SELECT id, username, email
	FROM users
	WHERE external_id = $1
	`

	err := s.pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting user by external id: %s", err)
	}

	return &user, nil
}

func (s *UserStore) GetById(ctx context.Context, id string) (*UserDisplay, error) {
	var user UserDisplay

	query := `
	SELECT u.id, u.username, u.email, gh.data ->> 'avatar_url' AS avatar_url
	FROM users u
	JOIN github_user_info gh ON u.id = gh.user_id
	WHERE u.id = $1
	`

	err := s.pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.AvatarURL)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting user by id: %s", err)
	}

	return &user, nil
}
