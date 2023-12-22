package graph

import (
	"fmt"

	"github.com/archspec/archspec-go/archspec/cpu"
	"github.com/converged-computing/jsongraph-go/jsongraph/v1/graph"
)

// getRootChild returns an edge that connects the root to a child node
func getRootChild(target *cpu.Microarchitecture, ids map[string]int) graph.Edge {
	fmt.Printf("Node %s is attached to the root\n", target.Name)
	m := getEdgeMetadata()
	return graph.Edge{Source: "0", Target: fmt.Sprintf("%d", ids[target.Name]), Metadata: m}
}
