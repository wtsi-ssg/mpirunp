package mpirun

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wtsi-ssg/mpirunp/port"
)

func TestMPIRun(t *testing.T) {
	Convey("Given a Wrapper made with an outDir and port range", t, func() {
		dir, err := ioutil.TempDir("", "mpiruntest")
		So(err, ShouldBeNil)
		defer os.RemoveAll(dir)

		outDir := filepath.Join(dir, "out")
		err = os.MkdirAll(outDir, 0770)
		So(err, ShouldBeNil)
		exists, err := pathExists(outDir)
		So(err, ShouldBeNil)
		So(exists, ShouldBeTrue)

		wrapper, err := New(outDir, 100, 134)
		So(err, ShouldBeNil)
		So(wrapper, ShouldNotBeNil)

		exists, err = pathExists(outDir)
		So(err, ShouldBeNil)
		So(exists, ShouldBeFalse)

		cmdString := "-output-filename " + outDir + " --mca oob_tcp_dynamic_ipv4_ports 100-116 --mca btl_tcp_port_min_v4 117 --mca btl_tcp_port_range_v4 17 foo bar"

		Convey("The correct mpirun command line is constructed", func() {
			So(wrapper.cmdArgs([]string{"foo", "bar"}), ShouldResemble,
				[]string{"-output-filename", outDir,
					"--mca", "oob_tcp_dynamic_ipv4_ports", fmt.Sprintf("%d-%d", 100, 116),
					"--mca", "btl_tcp_port_min_v4", "117", "--mca", "btl_tcp_port_range_v4", "17",
					"foo", "bar"})

			cmd := wrapper.Command("foo", "bar")
			So(cmd.String(), ShouldEndWith, "mpirun "+cmdString)
		})

		var outb, errb bytes.Buffer
		wrapper.Stdout = &outb
		wrapper.Stderr = &errb

		Convey("Execution fails if the output directory doesn't get created", func() {
			wrapper.exe = "echo"
			err = wrapper.Execute("foo", "bar")
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, ErrStuck)
			So(outb.String(), ShouldEqual, cmdString+"\n")
			So(errb.String(), ShouldBeBlank)
		})

		Convey("Execution succeeds if the output directory gets created", func() {
			wrapper.exe = "echo"
			err = os.MkdirAll(outDir, 0770)
			So(err, ShouldBeNil)

			err = wrapper.Execute("foo", "bar")
			So(err, ShouldBeNil)
			So(outb.String(), ShouldEqual, cmdString+"\n")
			So(errb.String(), ShouldBeBlank)
		})

		Convey("Execution succeeds with real mpirun", func() {
			ok := wrapper.CheckVersion()
			if !ok {
				SkipConvey("Skipping real mpirun test since Open MPI v4 isn't first in PATH", func() {})
			} else {
				checker, err := port.NewChecker("localhost")
				So(err, ShouldBeNil)
				min, max, err := checker.AvailableRange(3)
				So(err, ShouldBeNil)

				wrapper, err := New(outDir, min, max)
				So(err, ShouldBeNil)
				wrapper.wait = 1 * time.Second
				err = wrapper.Execute("sleep", "2")
				So(err, ShouldBeNil)
				So(outb.String(), ShouldBeBlank)
				So(errb.String(), ShouldBeBlank)
				exists, err = pathExists(outDir)
				So(err, ShouldBeNil)
				So(exists, ShouldBeTrue)

				Convey("But fails if the output dir is missing", func() {
					go func() {
						<-time.After(500 * time.Millisecond)
						err := os.RemoveAll(outDir)
						if err != nil {
							fmt.Printf("outdir removal failed: %s\n", err)
						}
					}()

					err = wrapper.Execute("sleep", "2")
					So(err, ShouldNotBeNil)
					So(err, ShouldEqual, ErrStuck)
				})
			}
		})
	})
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
