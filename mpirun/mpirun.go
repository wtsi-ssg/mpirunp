package mpirun

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const exe = "mpirun"
const exeV = "mpirun (Open MPI) 4"
const outArg = "-output-filename"
const mcaArg = "--mca"
const oobArg = "oob_tcp_dynamic_ipv4_ports"
const btlMinArg = "btl_tcp_port_min_v4"
const btlRangeArg = "btl_tcp_port_range_v4"
const outWait = 30 * time.Second

type constError string

func (err constError) Error() string {
	return string(err)
}

const (
	ErrStuck = constError("mpirun is non-responsive (failed to create output directory)")
)

// Wrapper is used to os/exec.Run mpirun, but wraps it with defining port ranges
// to use and checking if the output directory gets created, which it won't if
// it's just going to time out and kill itself.
type Wrapper struct {
	exe      string
	outDir   string
	oob      string
	btlMin   string
	btlRange string
	Stdout   io.Writer
	Stderr   io.Writer
	wait     time.Duration
}

// New creates a new Wrapper. outDir will be deleted immediately! When you
// Execute(), mpirun will use ports in the given range.
func New(outDir string, minPort, maxPort int) (*Wrapper, error) {
	err := os.RemoveAll(outDir)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	hosts := (maxPort - minPort - 1) / 2
	oobMax := minPort + hosts
	return &Wrapper{
		exe:      exe,
		outDir:   outDir,
		oob:      fmt.Sprintf("%d-%d", minPort, oobMax),
		btlMin:   strconv.Itoa(oobMax + 1),
		btlRange: strconv.Itoa(hosts + 1),
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
		wait:     outWait,
	}, nil
}

// Command generates an exec.Cmd with the correct output and port arguments,
// suffixed with the supplied args. You should probably just use Execute()
// instead.
func (w *Wrapper) Command(args ...string) *exec.Cmd {
	return exec.Command(w.exe, w.cmdArgs(args)...)
}

func (w *Wrapper) cmdArgs(args []string) []string {
	return append([]string{outArg, w.outDir,
		mcaArg, oobArg, w.oob,
		mcaArg, btlMinArg, w.btlMin,
		mcaArg, btlRangeArg, w.btlRange},
		args...)
}

// CheckVersion checks that openmpi is available, and is Open MPI version 4.
func (w *Wrapper) CheckVersion() bool {
	cmd := exec.Command(w.exe, "-V")
	out, err := cmd.Output()
	if err != nil || !strings.HasPrefix(string(out), exeV) {
		return false
	}
	return true
}

// Execute runs mpirun with the correct output and port arguments, plus the
// supplied args. It checks that the output directory gets created by mpirun
// after a little while, and if not kills it early.
func (w *Wrapper) Execute(args ...string) error {
	cmd := w.Command(args...)
	cmd.Stdout = w.Stdout
	cmd.Stderr = w.Stderr

	err := cmd.Start()
	if err != nil {
		return err
	}

	done := make(chan error, 1)
	go func() {
		err = cmd.Wait()
		done <- err
	}()

	timer := time.NewTimer(w.wait)

	select {
	case <-timer.C:
		if !w.outDirExists() {
			cmd.Process.Kill()
			timer.Stop()
			<-done

			return ErrStuck
		}
	case err = <-done:
		timer.Stop()
		if !w.outDirExists() {
			return ErrStuck
		}

		return err
	}

	return <-done
}

func (w *Wrapper) outDirExists() bool {
	_, err := os.Stat(w.outDir)
	return err == nil
}
