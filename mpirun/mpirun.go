package mpirun

import (
	"fmt"
	"os/exec"
	"strconv"
)

const exe = "mpirun"
const mcaArg = "--mca"
const oobArg = "oob_tcp_dynamic_ipv4_ports"
const btlMinArg = "btl_tcp_port_min_v4"
const btlRangeArg = "btl_tcp_port_range_v4"

type Wrapper struct {
	oob      string
	btlMin   string
	btlRange string
}

func New(minPort, maxPort int) *Wrapper {
	hosts := (maxPort - minPort - 1) / 2
	oobMax := minPort + hosts
	return &Wrapper{
		oob:      fmt.Sprintf("%d-%d", minPort, oobMax),
		btlMin:   strconv.Itoa(oobMax + 1),
		btlRange: strconv.Itoa(hosts + 1),
	}
}

func (w *Wrapper) Command(args ...string) *exec.Cmd {
	return exec.Command(exe, w.cmdArgs(args)...)
}

func (w *Wrapper) cmdArgs(args []string) []string {
	// --mca oob_tcp_dynamic_ipv4_ports 46100-46117 --mca btl_tcp_port_min_v4 46118 --mca btl_tcp_port_range_v4 17
	return append([]string{mcaArg, oobArg, w.oob, mcaArg, btlMinArg, w.btlMin, mcaArg, btlRangeArg, w.btlRange}, args...)
}
