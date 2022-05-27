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
		if visited[node] == NotVisited {
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
		if visited[to] == Entered {
			begin := len(path) - 1
			path = append(path, to)
			for path[begin] != to {
				begin -= 1
			}
			return path[begin:], true
		} else if visited[to] == NotVisited {
			if path, ok := g.isCyclic(to, visited, path); ok {
				return path, true
			}
		}
	}
	visited[node] = Exited
	return nil, false
}

func (g *Graph) TopSort() ([]string, bool) {
	visited := make(map[string]VisitedType)
	var topsort []string
	var ok bool
	for node := range g.nodes {
		if visited[node] == NotVisited {
			if topsort, ok = g.topSort(node, visited, topsort); !ok {
				return nil, false
			}
		}
	}
	return topsort, true
}

func (g *Graph) topSort(node string, visited map[string]VisitedType, topsort []string) ([]string, bool) {
	visited[node] = Entered
	for _, to := range g.adjacencyList[node] {
		if visited[to] == NotVisited {
			var ok bool
			if topsort, ok = g.topSort(to, visited, topsort); !ok {
				return nil, true
			}
		} else if visited[to] == Entered {
			return nil, false
		}
	}
	visited[node] = Exited
	topsort = append(topsort, node)
	return topsort, true
}
