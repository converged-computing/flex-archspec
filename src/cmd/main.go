package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/converged-computing/flex-archspec/src/graph"
)

func main() {
	fmt.Println("This is the flex archspec prototype")
	specFilePath := flag.String("spec", "", "JobSpec (yaml file) that defines machine architecture")
	matchPolicy := flag.String("policy", "first", "Match policy")
	saveFile := flag.String("file", "", "Save JGF to file (for debugging, etc.)")
	flag.Parse()

	specFile := *specFilePath

	// The JobSpec file is required
	if specFile == "" {
		flag.Usage()
		os.Exit(0)
	}

	// The JobSpec file and graphml must exist
	if _, err := os.Stat(specFile); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("%s does not exist\n", specFile)
		os.Exit(0)
	}

	// Create an ice cream graph, and match the spec to it.
	g := graph.FlexGraph{}
	g.Init(*matchPolicy, *saveFile)
	match, err := g.Match(specFile)
	if err != nil {
		fmt.Printf("Oh no! There was a problem with your machine matching: %x", err)
		return
	}
	match.Show()
}
