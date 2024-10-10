package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Node struct {
	Id         string
	Host       string
	Name       string
	User       string
	PemFile    string
	Port       int64
	Status     string
	AgentToken string
}

type NodeStore struct {
	pool *pgxpool.Pool
}

func (s *NodeStore) All(ctx context.Context) ([]Node, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id,
			host,
			name,
			username,
			status,
			convert_from(decode(pem_file, 'base64'), 'UTF8') as pem_file,
			agent_token
		FROM nodes
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []Node

	for rows.Next() {
		var node Node

		err := rows.Scan(&node.Id, &node.Host, &node.Name, &node.User, &node.Status, &node.PemFile, &node.AgentToken)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}

func (s *NodeStore) Create(ctx context.Context, node Node) (string, error) {
	var id string

	query := `
	INSERT INTO nodes (host, username, name, pem_file, port, agent_token)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`

	err := s.pool.QueryRow(ctx, query, node.Host, node.User, node.Name, node.PemFile, node.Port, node.AgentToken).Scan(&id)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *NodeStore) GetById(ctx context.Context, id string) (*Node, error) {
	var node Node

	query := `
	SELECT id, host, name, username, status, pem_file, agent_token
	FROM nodes
	WHERE id = $1
	`

	err := s.pool.QueryRow(ctx, query, id).Scan(&node.Id, &node.Host, &node.Name, &node.User, &node.Status, &node.PemFile, &node.AgentToken)

	if err != nil {
		return nil, err
	}

	return &node, nil
}

func (s *NodeStore) UpdateStatus(ctx context.Context, nodeID string, status string) error {
	query := `
	UPDATE nodes
	SET status = $1, updated_at = now()
	WHERE id = $2
	`

	_, err := s.pool.Exec(ctx, query, status, nodeID)

	if err != nil {
		return err
	}

	return nil
}
