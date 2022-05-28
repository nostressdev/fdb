package graph

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func FuzzGraph0(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string, cnt int) {
		names := strings.Split(s, "-")
		if len(names) >= cnt*2 && cnt > 0 {
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
			} else {
				assert.Nil(t, path)
			}
		}
	})
}

func FuzzGraph1(f *testing.F) {
	f.Fuzz(func(t *testing.T, cnt int, hasEdge []byte) {
		if len(hasEdge) == cnt*cnt {
			g := New()
			names := make([]string, 0)
			for i := 0; i < cnt; i++ {
				names = append(names, fmt.Sprintf("%d", i))
				g.AddNode(names[i])
			}
			edges := map[edge]bool{}
			for i := 0; i < cnt; i++ {
				for j := 0; j < cnt; j++ {
					if hasEdge[i*cnt+j] == 0 {
						edges[edge{names[i], names[j]}] = true
						g.AddEdge(names[i], names[j])
					}
				}
			}
			if ok, path := g.IsCyclic(); ok {
				for i := 0; i+1 < len(path); i += 1 {
					assert.True(t, edges[edge{path[i], path[i+1]}])
				}
			} else {
				assert.Nil(t, path)
			}
		}
	})
}
