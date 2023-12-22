# Flex Archspec

> Determine compatibility of your system with a container specification.

Context: The [OCI compatibility working group](https://github.com/opencontainers/wg-image-compatibility) is designing a new artifact to determine container image compatibility with a host environment of interest. We would want to be able to represent a specific compatibility artifact as a graph. Starting with [Archspec](https://github.com/archspec) is a logic route in that we can locate our micro-architecture of interest in the tree, and then determine if it is compatible with some container.

**Important** For current development, archspec json metadata requires a git submodule. I had to manually cd into the go directory where it was added and ensure it was cloned. E.g.,

```bash
go get github.com/archspec/archspec-go
# Your commit path may be different here
cd /home/vscode/go/pkg/mod/github.com/archspec/archspec-go\@v0.0.0-20231117085542-f806bb25b479/archspec
sudo git clone https://github.com/archspec/archspec-json json
```

This is reflected in [this issue](https://github.com/archspec/archspec-go/issues/13).

## Design

### Overview

The Flux Framework "flux-sched" or fluxion project provides modular bindings in different languages for intelligent,
graph-based scheduling. When we extend fluxion to a tool or project that warrants logic of this type, we call this a flex!
Thus, the project here demonstrates flex-archspec, or using fluxion to match some system request to what we know in the archspec graph. E.g.,

> Do you have an x86 system with this compiler option?

This is a simple use case that doesn't perfectly reflect the OCI container use case, but we need to start somewhere! For this very basic setup we are going to:

1. Load the machines into a JSON Graph (called JGF).
2. Try doing a query against system metadata

There will eventually be a third component - a container image specification, for which we need to include somewhere here. I am starting simple! 

Update: I have the graph loaded (and valid) and now need some help with understanding how it's traversed (e.g., and how to best query what I've represented, and if what I've represented is correct).


### Concepts

From the above, the following definitions might be useful.

 - **[Flux Framework](https://flux-framework.org)**: a modular framework for putting together a workload manager. It is traditionally for HPC, but components have been used in other places (e.g., here, Kubernetes, etc). It is analogous to Kubernetes in that it is modular and used for running batch workloads.
 - **[fluxion](fluxion)**: refers to [flux-framework/flux-sched](https://github.com/flux-framework/flux-sched) and is the scheduler component or module of Flux Framework. There are bindings in several languages, and specifically the Go bindings (server at [flux-framework/flux-k8s](https://github.com/flux-framework/flux-k8s)) assemble into the project "fluence."
 - **flex** is an out of tree tool, plugin, or similar that uses fluxion to flexibly schedule or match some kind of graph-based resources. This project is an example of a flex!

## Usage

### Build

This demonstrates how to build the bindings. You will need to be in the VSCode developer container environment, or produce the same
on your host. Note that we are currently relying on several WIP branches (or need/suggest changes to fluxion or Go bindings):

- We currently are using [this commit](https://github.com/researchapps/flux-sched/commit/0f33b17f6e792c14a262609d71f4ea5f32cb3ebb) that is a fork of [milroy's work](https://github.com/flux-framework/flux-sched/pull/1120) to ensure the module name matches what is added to go.mod.
- That branch also has added better error parsing as shown in [this issue](https://github.com/flux-framework/flux-sched/issues/1128)

When this is merged / the work is done, we will update to flux-framework/flux-sched. Below shows the make command that builds our final binary!

```bash
make
```
```console
# This needs to match the flux-sched install and latest commit, for now we are using a fork of milroy's branch
# that has a go.mod updated to match the org name
# go get -u github.com/researchapps/flux-sched/resource/reapi/bindings/go/src/fluxcli@86f5bb331342f2883b057920cf58e2c042aef881
go mod tidy
mkdir -p ./bin
GOOS=linux CGO_CFLAGS="-I/opt/flux-sched/resource/reapi/bindings/c" CGO_LDFLAGS="-L/usr/lib -L/opt/flux-sched/resource -lfluxion-resource -L/opt/flux-sched/resource/libjobspec -ljobspec_conv -L//opt/flux-sched/resource/reapi/bindings -lreapi_cli -lflux-idset -lstdc++ -lczmq -ljansson -lhwloc -lboost_system -lflux-hostlist -lboost_graph -lyaml-cpp" go build -ldflags '-w' -o bin/archspec src/cmd/main.go
```

The output is generated in bin:

```bash
$ ls bin/
archspec
```

### Run

Let's provide a 
You can provide your request for ice cream (e.g., icecream.yaml) and the description of the graph (in graphml). Note that we need shared libs on the path:

```bash
export LD_LIBRARY_PATH=/usr/lib:/opt/flux-sched/resource:/opt/flux-sched/resource/reapi/bindings:/opt/flux-sched/resource/libjobspec
```
```bash
./bin/archspec -spec ./examples/machine.yaml
```

It will save the JGF to a temporary json file (and remove it), but to save to one you can later inspect the graph:

```bash
./bin/archspec -spec ./examples/machine.yaml --file ./machine-graph.json
```
Note that I'm including the example [machine-graph.json](machine-graph.json) for inspection - I suspect there is a detail wrong about the structure (and how I'm querying) and will need to look closer at the details.

## License

HPCIC DevTools is distributed under the terms of the MIT license.
All new contributions must be made under this license.

See [LICENSE](https://github.com/converged-computing/cloud-select/blob/main/LICENSE),
[COPYRIGHT](https://github.com/converged-computing/cloud-select/blob/main/COPYRIGHT), and
[NOTICE](https://github.com/converged-computing/cloud-select/blob/main/NOTICE) for details.

SPDX-License-Identifier: (MIT)

LLNL-CODE- 842614
