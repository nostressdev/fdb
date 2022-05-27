package graph

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type edge struct {
	from string
	to   string
}

func TestGraph_IsCyclic(t *testing.T) {
	type fields struct {
		nodes         []string
		adjacencyList []edge
	}
	tests := []struct {
		name        string
		fields      fields
		hasCycle    bool
		path        []string
		repeatCount int
	}{
		{
			name: "no cycles",
			fields: fields{
				nodes: []string{"a", "b", "c", "d", "e"},
				adjacencyList: []edge{
					{"a", "b"},
					{"a", "c"},
					{"b", "d"},
					{"b", "e"},
					{"c", "b"},
					{"e", "d"},
				},
			},
			hasCycle:    false,
			path:        nil,
			repeatCount: 1000,
		},
		{
			name: "cycle",
			fields: fields{
				nodes: []string{"a", "b", "c", "d", "e"},
				adjacencyList: []edge{
					{"a", "b"},
					{"a", "c"},
					{"b", "d"},
					{"c", "b"},
					{"d", "a"},
				},
			},
			hasCycle: true,
			path: []string{
				"a", "b", "d",
			},
			repeatCount: 100,
		},
		{
			name: "cycle",
			fields: fields{
				nodes: []string{"a", "b"},
				adjacencyList: []edge{
					{"a", "b"},
					{"b", "a"},
				},
			},
			hasCycle: true,
			path: []string{
				"a", "b",
			},
			repeatCount: 100,
		},

		{
			name: "no cycle",
			fields: fields{
				nodes: []string{"a", "b", "c", "d", "e"},
				adjacencyList: []edge{
					{"a", "b"},
					{"b", "c"},
					{"a", "c"},
				},
			},
			hasCycle:    false,
			path:        nil,
			repeatCount: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for cnt := 0; cnt < tt.repeatCount; cnt++ {
				g := New()
				for _, node := range tt.fields.nodes {
					g.AddNode(node)
				}
				for _, edge := range tt.fields.adjacencyList {
					g.AddEdge(edge.from, edge.to)
				}
				hasCycle, path := g.IsCyclic()
				require.Equal(t, hasCycle, tt.hasCycle)
				if hasCycle {
					require.True(t, len(path) > 1, "cycle path cannot be 1 in len")
					path = path[1:]
					if len(path) == 1 {
						assert.True(t, reflect.DeepEqual(path, tt.path))
					} else {
						test := false
						for index := 0; !test && index < len(path); index++ {
							path = append(path[1:], path[:1]...)
							test = test || reflect.DeepEqual(path, tt.path)
						}
						assert.True(t, test)
					}
				} else {
					assert.Nil(t, path, "np cycle path must be empty")
				}
			}
		})
	}
}

func TestGraph_TopSort(t *testing.T) {
	type fields struct {
		nodes         []string
		adjacencyList []edge
	}
	tests := []struct {
		name     string
		fields   fields
		hasCycle bool
	}{
		{
			name: "one variant",
			fields: fields{
				nodes: []string{"a", "b", "c"},
				adjacencyList: []edge{
					{"a", "b"},
					{"b", "c"},
					{"a", "c"},
				},
			},
			hasCycle: false,
		},
		{
			name: "many variants",
			fields: fields{
				nodes: []string{"a", "b", "c", "d", "e"},
				adjacencyList: []edge{
					{"a", "b"},
					{"b", "c"},
					{"a", "c"},
					{"d", "e"},
				},
			},
			hasCycle: false,
		},
		{
			name: "cycle",
			fields: fields{
				nodes: []string{"a", "b", "c", "d", "e"},
				adjacencyList: []edge{
					{"a", "b"},
					{"b", "c"},
					{"c", "a"},
				},
			},
			hasCycle: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New()
			for _, node := range tt.fields.nodes {
				g.AddNode(node)
			}
			for _, edge := range tt.fields.adjacencyList {
				g.AddEdge(edge.from, edge.to)
			}
			topsort, ok := g.TopSort()
			assert.Equal(t, !ok, tt.hasCycle, "cycle status mismatched")
			if !tt.hasCycle {
				nodesSet := make(map[string]bool)
				for _, node := range topsort {
					for _, to := range g.adjacencyList[node] {
						require.Equal(t, true, nodesSet[to], "graph is not topsorted: %s -> %s", node, to)
					}
					nodesSet[node] = true
				}
			}
		})
	}
}
