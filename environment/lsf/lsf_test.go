package lsf

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLSFHosts(t *testing.T) {
	Convey("You can get LSB_HOSTS", t, func() {
		hosts := Hosts()
		So(len(hosts), ShouldEqual, 0)

		os.Setenv(hostsEnv, "foo bar")
		hosts = Hosts()
		So(len(hosts), ShouldEqual, 2)
		So(hosts, ShouldResemble, []string{"foo", "bar"})
	})
}
