package graph

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func FuzzGraphFindCycle(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string, cnt int) {
		names := strings.Split(s, "-")
		if len(names) >= cnt && cnt > 0 {
			if (len(names)-cnt)%2 == 1 {
				cnt--
			}
			g := New()
			for i := 0; i < cnt; i++ {
				g.AddNode(names[i])
			}
			edges := map[edge]bool{}
			for i := cnt; i < len(names); i += 2 {
				edges[edge{names[i], names[i+1]}] = true
				g.AddEdge(names[i], names[i+1])
			}
			if ok, path := g.IsCyclic(); ok {
				for i := 0; i+1 < len(path); i += 1 {
					assert.True(t, edges[edge{path[i], path[i+1]}])
				}
			}

		}
	})
}
