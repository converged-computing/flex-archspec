package archspec

import "fmt"

// A MachineRequest asks for a machine
type MachineRequest struct {
	Number uint64
	Spec   string
}

func (i *MachineRequest) Satisfied() bool {
	return i.Spec != ""
}

// Show the customer their final request
func (i *MachineRequest) Show() {
	if i.Satisfied() {
		fmt.Printf("\n😍️ Your Machine Arch is Satisfiable!\n")
		fmt.Printf("Number: %d\n", i.Number)
		fmt.Printf("  Spec:\n%s", i.Spec)
	} else {
		fmt.Printf("\n😭️ Oh no, we could not satisfy your request!\n")
	}
}
