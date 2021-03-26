package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wtsi-ssg/mpirunp/environment/lsf"
	"github.com/wtsi-ssg/mpirunp/mpirun"
	"github.com/wtsi-ssg/mpirunp/port"
)

func main() {
	var outPath string
	flag.StringVar(&outPath, "output-filename", "", "[required] Redirect output from application processes into filename/rank.out (will be treated as a directory that will be deleted and then created at start up)")
	flag.Parse()
	if outPath == "" {
		fmt.Printf("-output-filename is requred\n")
		os.Exit(1)
	}

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

	wrapper, err := mpirun.New(outPath, min, max)
	if err != nil {
		fmt.Printf("Failure preparing the output directory: %s\n", err)
		os.Exit(1)
	}

	if ok := wrapper.CheckVersion(); !ok {
		fmt.Printf("Only Open MPI v4 is supported\n")
		os.Exit(1)
	}

	cmd := wrapper.Command(flag.Args()...)
	fmt.Printf("Will run: %s\n", cmd.String())

	err = wrapper.Execute(flag.Args()...)
	if err != nil {
		fmt.Printf("Execution of mpirun failed: %s\n", err)
		os.Exit(1)
	}
}
