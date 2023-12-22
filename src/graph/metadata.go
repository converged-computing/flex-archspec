package graph

import (
	"github.com/archspec/archspec-go/archspec/cpu"
	"github.com/converged-computing/jsongraph-go/jsongraph/metadata"
	// TODO update back to flux-sched when merged
)

/*

Desired steps:

1. Load the machines into a JSON Graph (called JGF).
2. Try doing a query against system metadata

*/

const (
	microArchType = "microarchitecture"
)

var (
	rootNode   = "machine"
	rootVendor = "generic"
)

// getEdgeMetadata returns default edge metadata.
// We assume an "in" relationship of a node being in (a child of) a parent
func getEdgeMetadata() metadata.Metadata {
	// Make up some containment metadata
	// "metadata": {
	//	"name": {
	//	  "containment": "in"
	//	}
	//  }
	m := metadata.Metadata{}
	nameKey := map[string]string{"containment": "in"}
	m.AddElement("name", nameKey)
	return m
}

// getTargetMetadata starts with default metadata and adds on target specific attributes
func getTargetMetadata(target *cpu.Microarchitecture, ids map[string]int) metadata.Metadata {

	m := getDefaultMetadata()
	counter := ids[target.Name]
	path := getTargetPath(target, ids)
	m.AddElement("vendor", target.Vendor)
	m.AddElement("generation", target.Generation)
	m.AddElement("name", target.Name)
	m.AddElement("uniq_id", counter)
	m.AddElement("id", counter)
	m.AddElement("paths", map[string]string{"containment": path})
	m.AddElement("type", target.Name)

	// Features are like metadata keys?
	// Treat each feature like a boolean (yes/no)
	target.Features.Each(func(item string) bool {
		m.AddElement(item, "yes")
		return true
	})
	return m
}

// getDefaultMetadata ensures required fields (that aren't specific to a target) are added
func getDefaultMetadata() metadata.Metadata {

	m := metadata.Metadata{}

	// These are required metadata fields
	// See https://github.com/flux-framework/flux-sched/blob/745e3e097fe1368e53fcf46b0a2c2cd7b95ad369/resource/readers/resource_reader_jgf.cpp#L383-L389
	m.AddElement("type", microArchType)
	m.AddElement("basename", microArchType)
	m.AddElement("rank", -1)
	m.AddElement("status", -1)
	m.AddElement("exclusive", false)
	m.AddElement("unit", "")
	m.AddElement("size", 1)
	return m
}
