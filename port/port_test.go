package port

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPort(t *testing.T) {
	Convey("Given a Checker", t, func() {
		checker, err := NewChecker("localhost")
		So(err, ShouldBeNil)
		So(checker, ShouldNotBeNil)

		Convey("You can get an available port number", func() {
			port, err := checker.availablePort()
			So(err, ShouldBeNil)
			So(port, ShouldBeBetweenOrEqual, 1, maxPort)
			So(len(checker.ports), ShouldEqual, 1)
			So(checker.ports[port], ShouldBeTrue)

			err = checker.release()
			So(err, ShouldBeNil)
		})

		Convey("portsAfter works", func() {
			after := checker.portsAfter(10)
			So(len(after), ShouldEqual, 0)

			checker.ports[9] = true
			checker.ports[12] = true
			checker.ports[13] = true
			checker.ports[15] = true

			after = checker.portsAfter(10)
			So(len(after), ShouldEqual, 0)

			checker.ports[11] = true
			after = checker.portsAfter(10)
			So(len(after), ShouldEqual, 3)
			So(after, ShouldResemble, []int{11, 12, 13})
		})

		Convey("portsBefore works", func() {
			after := checker.portsBefore(10)
			So(len(after), ShouldEqual, 0)

			checker.ports[11] = true
			checker.ports[8] = true
			checker.ports[7] = true
			checker.ports[5] = true

			after = checker.portsBefore(10)
			So(len(after), ShouldEqual, 0)

			checker.ports[9] = true
			after = checker.portsBefore(10)
			So(len(after), ShouldEqual, 3)
			So(after, ShouldResemble, []int{7, 8, 9})
		})

		Convey("checkRange returns nothing with no available ports", func() {
			set, has := checker.checkRange(10, 4)
			So(has, ShouldBeFalse)
			So(len(set), ShouldEqual, 0)

			Convey("but returns ports above starting point", func() {
				checker.ports[9] = true
				checker.ports[11] = true
				checker.ports[12] = true
				checker.ports[13] = true
				checker.ports[14] = true

				set, has := checker.checkRange(10, 4)
				So(has, ShouldBeTrue)
				So(len(set), ShouldEqual, 4)
				So(set, ShouldResemble, []int{10, 11, 12, 13})
			})

			Convey("but returns ports below starting point", func() {
				checker.ports[11] = true
				checker.ports[9] = true
				checker.ports[8] = true
				checker.ports[7] = true
				checker.ports[6] = true

				set, has := checker.checkRange(10, 4)
				So(has, ShouldBeTrue)
				So(len(set), ShouldEqual, 4)
				So(set, ShouldResemble, []int{7, 8, 9, 10})
			})

			Convey("but returns ports below and above starting point", func() {
				checker.ports[8] = true
				checker.ports[9] = true
				checker.ports[11] = true
				checker.ports[12] = true

				set, has := checker.checkRange(10, 4)
				So(has, ShouldBeTrue)
				So(len(set), ShouldEqual, 4)
				So(set, ShouldResemble, []int{8, 9, 10, 11})
			})

			Convey("and returns nothing with non-contiguous available ports", func() {
				checker.ports[7] = true
				checker.ports[8] = true
				checker.ports[12] = true
				checker.ports[13] = true

				set, has := checker.checkRange(10, 4)
				So(has, ShouldBeFalse)
				So(len(set), ShouldEqual, 0)
			})
		})

		Convey("You can get a range of available ports", func() {
			ports, err := checker.AvailableRange(2)
			So(err, ShouldBeNil)
			So(len(ports), ShouldEqual, 2)
			So(ports[0], ShouldBeBetweenOrEqual, 1, maxPort)
			So(ports[1], ShouldEqual, ports[0]+1)

			ports, err = checker.AvailableRange(67)
			So(err, ShouldBeNil)
			So(len(ports), ShouldEqual, 67)
			So(ports[0], ShouldBeBetweenOrEqual, 1, maxPort)
			So(ports[1], ShouldEqual, ports[0]+1)
			So(ports[66], ShouldEqual, ports[0]+66)
		})
	})
}