package mpirun

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMPIRun(t *testing.T) {
	Convey("You can get an mpirun command using given ports", t, func() {
		wrapper := New(100, 134)
		So(wrapper, ShouldNotBeNil)

		So(wrapper.cmdArgs([]string{"foo", "bar"}), ShouldResemble, []string{"--mca", "oob_tcp_dynamic_ipv4_ports", fmt.Sprintf("%d-%d", 100, 116), "--mca", "btl_tcp_port_min_v4", "117", "--mca", "btl_tcp_port_range_v4", "17", "foo", "bar"})

		cmd := wrapper.Command("foo", "bar")
		So(cmd.String(), ShouldEndWith, "mpirun --mca oob_tcp_dynamic_ipv4_ports 100-116 --mca btl_tcp_port_min_v4 117 --mca btl_tcp_port_range_v4 17 foo bar")
	})
}
