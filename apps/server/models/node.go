package models

import (
	"database/sql"
)

type Node struct {
	Id      string `json:"id"`
	Host    string `json:"host"`
	Name    string `json:"name"`
	User    string `json:"user"`
	PemFile string `json:"pemFile"`
}

type NodeModel struct {
	DB *sql.DB
}

func (m NodeModel) All() ([]Node, error) {
	rows, err := m.DB.Query("SELECT id, host, name, user FROM nodes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []Node

	for rows.Next() {
		var node Node

		err := rows.Scan(&node.Id, &node.Host, &node.Name, &node.User)
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

func (m NodeModel) AddNode(node Node) error {
	query := `
	INSERT INTO nodes (host, username, name, pem_file)
	VALUES ($1, $2, $3, $4)
	`

	_, err := m.DB.Exec(query, node.Host, node.User, node.Name, node.PemFile)

	if err != nil {
		return err
	}

	return nil
}
