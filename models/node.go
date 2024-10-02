package models

import (
	"database/sql"
)

type Node struct {
	Id      string
	Host    string
	Name    string
	User    string
	PemFile string
	Port    int64
	Status  string
}

type NodeModel struct {
	DB *sql.DB
}

func (m NodeModel) All() ([]Node, error) {
	rows, err := m.DB.Query("SELECT id, host, name, user, status FROM nodes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []Node

	for rows.Next() {
		var node Node

		err := rows.Scan(&node.Id, &node.Host, &node.Name, &node.User, &node.Status)
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

func (m NodeModel) AddNode(node Node) (string, error) {
	var id string

	query := `
	INSERT INTO nodes (host, username, name, pem_file, port)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`

	err := m.DB.QueryRow(query, node.Host, node.User, node.Name, node.PemFile, node.Port).Scan(&id)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (m NodeModel) UpdateNodeStatus(nodeID string, status string) error {
	query := `
	UPDATE nodes
	SET status = $1, updated_at = now()
	WHERE id = $2
	`

	_, err := m.DB.Exec(query, status, nodeID)

	if err != nil {
		return err
	}

	return nil
}
