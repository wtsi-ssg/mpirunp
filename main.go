package main

import (
	"fmt"
	"os"

	"github.com/wtsi-ssg/mpirunp/environment/lsf"
	"github.com/wtsi-ssg/mpirunp/mpirun"
	"github.com/wtsi-ssg/mpirunp/port"
)

func main() {
	hosts := lsf.Hosts()
	fmt.Printf("Working with %d hosts.\n", len(hosts))
	portsNeeded := (len(hosts) * 2) + 2

	checker, err := port.NewChecker("localhost")
	if err != nil {
		fmt.Printf("Making a port checker failed: %s\n", err)
		os.Exit(1)
	}

	min, max, err := checker.AvailableRange(portsNeeded)
	if err != nil {
		fmt.Printf("Getting a range of %d contiguous ports failed: %s\n", portsNeeded, err)
		os.Exit(1)
	}

	fmt.Printf("Ports %d..%d are free right now.\n", min, max)

	wrapper := mpirun.New(min, max)
	cmd := wrapper.Command(os.Args[1:]...)
	fmt.Printf("Will run: %s\n", cmd.String())

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Execution of mpirun failed: %s\n", err)
		os.Exit(1)
	}
}
