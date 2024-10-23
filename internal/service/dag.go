package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type DAG struct {
	nodes map[uuid.UUID][]uuid.UUID
	indeg map[uuid.UUID]int
}

type Edge struct {
	Source uuid.UUID `json:"source"`
	Dest   uuid.UUID `json:"dest"`
}

func NewDAG() *DAG {
	return &DAG{
		nodes: make(map[uuid.UUID][]uuid.UUID),
		indeg: make(map[uuid.UUID]int),
	}
}

func (g *DAG) AddNode(node uuid.UUID) {
	if _, exists := g.nodes[node]; !exists {
		g.nodes[node] = []uuid.UUID{}
		g.indeg[node] = 0
	}
}

func (g *DAG) AddEdge(node1, node2 uuid.UUID) error {
	if _, exists := g.nodes[node1]; !exists {
		return errors.New("node1 does not exist")
	}
	if _, exists := g.nodes[node2]; !exists {
		return errors.New("node2 does not exist")
	}

	for _, neighbor := range g.nodes[node1] {
		if neighbor == node2 {
			return nil
		}
	}

	g.nodes[node1] = append(g.nodes[node1], node2)
	g.indeg[node2]++
	return nil
}

func (g *DAG) GetExecutionGroups() ([][]uuid.UUID, error) {
	var result [][]uuid.UUID

	queue := []uuid.UUID{}
	for node, deg := range g.indeg {
		if deg == 0 {
			queue = append(queue, node)
		}
	}

	for len(queue) > 0 {
		var nextQueue []uuid.UUID
		group := []uuid.UUID{}

		for _, node := range queue {
			group = append(group, node)

			for _, neighbor := range g.nodes[node] {
				g.indeg[neighbor]--
				if g.indeg[neighbor] == 0 {
					nextQueue = append(nextQueue, neighbor)
				}
			}
		}

		result = append(result, group)
		queue = nextQueue
	}

	for _, deg := range g.indeg {
		if deg > 0 {
			return nil, errors.New("cycle detected in DAG")
		}
	}

	return result, nil
}

func (g *DAG) GetEdges() []Edge {
	var edges []Edge
	for source, neighbors := range g.nodes {
		for _, dest := range neighbors {
			edges = append(edges, Edge{Source: source, Dest: dest})
		}
	}
	return edges
}

func (g *DAG) PrettyPrint() {
	for node, neighbors := range g.nodes {
		fmt.Printf("Node %s -> ", node)
		if len(neighbors) == 0 {
			fmt.Println("[]")
		} else {
			for i, neighbor := range neighbors {
				if i == len(neighbors)-1 {
					fmt.Printf("%s\n", neighbor)
				} else {
					fmt.Printf("%s, ", neighbor)
				}
			}
		}
	}
}
