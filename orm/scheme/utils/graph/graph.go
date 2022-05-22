package graph

type Graph struct {
	nodes         map[string]bool
	adjucencyList map[string][]string
}

func New() *Graph {
	return &Graph{
		nodes:         make(map[string]bool),
		adjucencyList: make(map[string][]string),
	}
}

func (g *Graph) AddNode(node string) {
	g.nodes[node] = true
}

func (g *Graph) AddEdge(from, to string) {
	g.adjucencyList[from] = append(g.adjucencyList[from], to)
}

func (g *Graph) IsCyclic() (bool, []string) {
	visited := make(map[string]int)
	path := make([]string, 0)
	for node := range g.nodes {
		if g.isCyclic(node, visited, path) {
			return true, path
		}
	}
	return false, nil
}

func (g *Graph) isCyclic(node string, visited map[string]int, path []string) bool {
	visited[node] = 1
	path = append(path, node)
	for _, to := range g.adjucencyList[node] {
		if visited[to] == 1 {
			return true
		} else if visited[to] == 0 {
			if g.isCyclic(to, visited, path) {
				return true
			}
		}
	}
	visited[node] = 2
	path = path[:len(path)-1]
	return false
}
