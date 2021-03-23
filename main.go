package main

import (
	"fmt"
	"os"

	"github.com/wtsi-ssg/mpirunp/port"
)

func main() {
	portsNeeded := 67

	checker, err := port.NewChecker("localhost")
	if err != nil {
		fmt.Printf("Making a port checker failed: %s\n", err)
		os.Exit(1)
	}

	ports, err := checker.AvailableRange(portsNeeded)
	if err != nil {
		fmt.Printf("Getting a range of %d contiguous ports failed: %s\n", portsNeeded, err)
		os.Exit(1)
	}

	fmt.Printf("got ports %v\n", ports)
}
