package graph

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/archspec/archspec-go/archspec/cpu"
	"github.com/converged-computing/flex-archspec/src/archspec"

	"github.com/converged-computing/jsongraph-go/jsongraph/v1/graph"

	// TODO update back to flux-sched when merged
	"github.com/researchapps/flux-sched/resource/reapi/bindings/go/src/fluxcli"

	"fmt"
)

/*

Desired steps:

1. Load the machines into a JSON Graph (called JGF).
2. Try doing a query against system metadata

*/

type FlexGraph struct {
	cli *fluxcli.ReapiClient
}

// Helper functions to parse a target
func getTargetPath(target *cpu.Microarchitecture, ids map[string]int) string {

	// generate the paths string from the parents
	path := target.Name
	for _, parent := range target.Parents {
		uid := ids[parent.Name]
		path = fmt.Sprintf("%s%d/%s", parent.Name, uid, path)
	}
	path = fmt.Sprintf("/%s/%s", rootNode, path)
	return path
}

// Init a new FlexGraph from a graphml filename
func (f *FlexGraph) Init(matchPolicy string, saveFile string) error {

	// 1. instantiate fluxion
	f.cli = fluxcli.NewReapiClient()
	fmt.Printf("Created flex resource graph %s\n", f.cli)

	// prepare a graph to load targets into
	g := graph.NewGraph()

	// Save target ids for later
	ids := map[string]int{}
	counter := 1

	// Generate UIDs first. These seem arbitrary
	for _, target := range cpu.TARGETS {
		ids[target.Name] = counter
		counter += 1
	}

	// Create an abstract root node
	root := generateRoot()
	g.Graph.Nodes = append(g.Graph.Nodes, root)

	// Show the targets (save ids as we go)
	for targetName, target := range cpu.TARGETS {

		counter := ids[target.Name]

		fmt.Printf("Adding node for %s\n", target.Name)
		fmt.Printf("    Generation: %d\n", target.Generation)
		fmt.Printf("        Vendor: %s\n", target.Vendor)

		// If we don't have parents, it's attached to the root node
		if len(target.Parents) == 0 {
			edge := getRootChild(&target, ids)
			g.Graph.Edges = append(g.Graph.Edges, edge)
		}

		// Generate metadata for the node
		m := getTargetMetadata(&target, ids)

		// Each target is a new Node
		uid := fmt.Sprintf("%d", counter)
		node := graph.Node{
			Label:    &targetName,
			Id:       uid,
			Metadata: m,
		}
		counter += 1
		g.Graph.Nodes = append(g.Graph.Nodes, node)

		// The edges are the parents
		for _, parent := range target.Parents {

			m = getEdgeMetadata()
			parentId := fmt.Sprintf("%d", ids[parent.Name])
			targetId := fmt.Sprintf("%d", ids[target.Name])

			// I assume no id means a root node
			if parentId == "" {
				parentId = rootNode
			}
			edge := graph.Edge{Source: parentId, Target: targetId, Metadata: m}
			g.Graph.Edges = append(g.Graph.Edges, edge)
		}
	}

	// Set match policy to default (first) if not defined.
	// In practice this should not happen - the cmd/main.go sets a default.
	if matchPolicy == "" {
		matchPolicy = "first"
	}

	// Serialize the struct to string
	conf, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return err
	}

	if saveFile == "" {
		jsonFile, err := os.CreateTemp("", "machine-*.json") // in Go version older than 1.17 you can use ioutil.TempFile
		if err != nil {
			fmt.Printf("Error creating temporary json file: %x", err)
			return err
		}
		defer jsonFile.Close()
		defer os.Remove(jsonFile.Name())
		saveFile = jsonFile.Name()
	}

	// Write to file!
	err = os.WriteFile(saveFile, conf, os.ModePerm)
	if err != nil {
		fmt.Printf("Error writing json to file: %x", err)
		return err
	}

	// Alert the user to all the chosen parameters
	// Note that "grug" == "graphml" but probably nobody knows what grug means
	fmt.Printf(" Match policy: %s\n", matchPolicy)
	fmt.Println(" Load format: JSON Graph Format (JGF)")

	// 2. Create the context, the default format is JGF
	// 3. Remainder of defaults should work out of the box
	// Note that the options get passed as a json string to here:
	// https://github.com/flux-framework/flux-sched/blob/master/resource/reapi/bindings/c%2B%2B/reapi_cli_impl.hpp#L412
	opts := `{"matcher_policy": "%s", "load_file": "%s", "load_format": "jgf", "match_format": "jgf"}`
	p := fmt.Sprintf(opts, matchPolicy, saveFile)

	// 4. Then pass in a jobspec... err, ice cream request :)
	err = f.cli.InitContext(string(conf), p)
	if err != nil {
		fmt.Printf("Error creating context: %s", err)
		return err
	}
	fmt.Printf("\n✨️ Init context complete!\n")
	return nil
}

// Order is akin to doing a Satisfies, but right now it's a MatchAllocate
// The result of an order is a traversal of the graph that could satisfy the ice cream request
func (f *FlexGraph) Match(specFile string) (archspec.MachineRequest, error) {
	fmt.Printf("💻️   Request: %s\n", specFile)

	// Prepare the ice cream request!
	request := archspec.MachineRequest{}

	spec, err := os.ReadFile(specFile)
	if err != nil {
		return request, errors.New("Error reading jobspec")
	}

	// TODO this could be f.cli.Satisfies
	// Note that number originally was a jobid (it's now a number for the ice cream in the shop)
	// Note that recipe was originally "allocated"
	_, machine, _, _, number, err := f.cli.MatchAllocate(false, string(spec))
	if err != nil {
		return request, err
	}

	// Populate the ice cream request for the customer
	request.Spec = machine
	request.Number = number
	return request, nil
}
