package graph

import (
	"reflect"
	"testing"
)

func TestGraph_IsCyclic(t *testing.T) {
	type fields struct {
		nodes         map[string]bool
		adjucencyList map[string][]string
	}
	tests := []struct {
		name   string
		fields fields
		hasCycle   bool
		path  []string
	}{
		{
			name: "no cycles",
			fields: fields{
				nodes: map[string]bool{
					"a": true,
					"b": true,
					"c": true,
					"d": true,
					"e": true,
				},
				adjucencyList: map[string][]string{
					"a": {"b", "c"},
					"b": {"d", "e"},
					"c": {"b"},
					"e": {"d"},
				},
			},
			hasCycle: false,
			path: nil,
		},
		{
			name: "cycle",
			fields: fields{
				nodes: map[string]bool{
					"a": true,
					"b": true,
					"c": true,
					"d": true,
				},
				adjucencyList: map[string][]string{
					"a": {"b", "c"},
					"b": {"d"},
					"c": {"b"},
					"d": {"a"},
				},
			},
			hasCycle: true,
			path: []string{
				"a", "b", "d", "a",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Graph{
				nodes:         tt.fields.nodes,
				adjucencyList: tt.fields.adjucencyList,
			}
			hasCycle, path := g.IsCyclic()
			if hasCycle != tt.hasCycle {
				t.Errorf("Graph.IsCyclic() got = %v, want %v", hasCycle, tt.hasCycle)
			}
			if !reflect.DeepEqual(path, tt.path) {
				t.Errorf("Graph.IsCyclic() got1 = %v, want %v", path, tt.path)
			}
		})
	}
}
