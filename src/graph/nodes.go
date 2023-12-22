package graph

import (
	"fmt"

	"github.com/converged-computing/jsongraph-go/jsongraph/v1/graph"
)

// generateRoot generates an abstract root for the graph
func generateRoot() graph.Node {

	// Generate metadata for the node
	m := getDefaultMetadata()
	m.AddElement("vendor", rootVendor)
	m.AddElement("name", rootNode)
	m.AddElement("uniq_id", 0)
	m.AddElement("id", 0)
	m.AddElement("paths", map[string]string{"containment": fmt.Sprintf("/%s", rootNode)})

	return graph.Node{
		Label:    &rootNode,
		Id:       "0",
		Metadata: m,
	}

}
