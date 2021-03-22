package main

import (
	"fmt"
	"time"

	"github.com/wtsi-ssg/mpirunp/portscan"
	"github.com/wtsi-ssg/mpirunp/ulimit"
	"golang.org/x/sync/semaphore"
)

func main() {
	limiter := semaphore.NewWeighted(ulimit.Ulimit())
	timeout := 500 * time.Millisecond

	ps := portscan.New("127.0.0.1", limiter, timeout)
	available := ps.AvailablePorts(1, 65535)
	fmt.Printf("%d available ports\n", len(available))

	free, err := portscan.GetFreePorts(1000)
	if err != nil {
		fmt.Printf("Free failed: %s\n", err)
	} else {
		fmt.Printf("%d free ports:\n%v\n", len(free), free)
	}
}
