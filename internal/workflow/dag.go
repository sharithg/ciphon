package workflow

import "errors"

type DAG struct {
	nodes   map[string][]string
	visited map[string]bool
	stack   []string
}

func NewDAG() *DAG {
	return &DAG{
		nodes:   make(map[string][]string),
		visited: make(map[string]bool),
	}
}

func (g *DAG) AddNode(node string) {
	if _, exists := g.nodes[node]; !exists {
		g.nodes[node] = []string{}
	}
}

func (g *DAG) AddEdge(node1, node2 string) error {
	if _, exists := g.nodes[node1]; !exists {
		return errors.New("node1 does not exist")
	}
	if _, exists := g.nodes[node2]; !exists {
		return errors.New("node2 does not exist")
	}

	g.nodes[node1] = append(g.nodes[node1], node2)
	return nil
}

func (g *DAG) TopologicalSort() ([]string, error) {
	g.visited = make(map[string]bool)
	g.stack = []string{}

	recStack := make(map[string]bool)

	for node := range g.nodes {
		if !g.visited[node] {
			if err := g.dfs(node, recStack); err != nil {
				return nil, err
			}
		}
	}

	result := make([]string, len(g.stack))
	for i, v := range g.stack {
		result[len(g.stack)-1-i] = v
	}
	return result, nil
}

func (g *DAG) dfs(node string, recStack map[string]bool) error {
	g.visited[node] = true
	recStack[node] = true

	for _, neighbor := range g.nodes[node] {
		if !g.visited[neighbor] {
			if err := g.dfs(neighbor, recStack); err != nil {
				return err
			}
		} else if recStack[neighbor] {
			return errors.New("cycle detected in DAG")
		}
	}

	recStack[node] = false
	g.stack = append(g.stack, node)
	return nil
}
