package graph

type VisitedType int

const (
	NotVisited = VisitedType(iota)
	Entered
	Exited
)

type Graph struct {
	nodes         map[string]bool
	adjacencyList map[string][]string
}

func New() *Graph {
	return &Graph{
		nodes:         make(map[string]bool),
		adjacencyList: make(map[string][]string),
	}
}

func (g *Graph) AddNode(node string) {
	g.nodes[node] = true
}

func (g *Graph) AddEdge(from, to string) {
	g.adjacencyList[from] = append(g.adjacencyList[from], to)
}

func (g *Graph) IsCyclic() (bool, []string) {
	visited := make(map[string]VisitedType)
	for node := range g.nodes {
		if visited[node] == 0 {
			if path, ok := g.isCyclic(node, visited, nil); ok {
				return true, path
			}
		}
	}
	return false, nil
}

func (g *Graph) isCyclic(node string, visited map[string]VisitedType, path []string) ([]string, bool) {
	visited[node] = Entered
	path = append(path, node)
	for _, to := range g.adjacencyList[node] {
		if visited[to] == 1 {
			begin := len(path) - 1
			path = append(path, to)
			for path[begin] != to {
				begin -= 1
			}
			path = path[begin:]
			return path, true
		} else if visited[to] == 0 {
			if path, ok := g.isCyclic(to, visited, path); ok {
				return path, true
			}
		}
	}
	visited[node] = Exited
	return nil, false
}
