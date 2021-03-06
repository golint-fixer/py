package py

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/py.v0/mainthread"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"testing"
	"time"
)

func TestConvertPyObject2Go(t *testing.T) {
	Convey("Given an initialized python py2go test module", t, func() {
		mainthread.AppendSysPath("")

		mdl, err := LoadModule("_test_py2go")
		So(err, ShouldBeNil)
		So(mdl, ShouldNotBeNil)

		returnTypes := []struct {
			typeName string
			expected data.Value
		}{
			{"true", data.Bool(true)},
			{"false", data.Bool(false)},
			{"int", data.Int(123)},
			{"float", data.Float(1.0)},
			{"string", data.String("ABC")},
			{"unicode", data.String("hello")},
			{"bytearray", data.Blob([]byte("abcdefg"))},
			{"map", data.Map{"key1": data.Int(123), "key2": data.String("str")}},
			{"nested_map", data.Map{"key1": data.Map{"key2": data.Int(123)}}},
			{"array", data.Array{data.Int(1), data.Int(2), data.Map{"key": data.Int(3)}}},
			{"none", data.Null{}},
			{"timestamp", data.Timestamp(time.Date(2015, time.May, 1, 14, 27, 0, 500*int(time.Millisecond), time.UTC))},
			{"timestamp_with_tz", data.Timestamp(time.Date(2015, time.May, 1, 5, 24, 0, 500*int(time.Millisecond), time.UTC))},
			{"onetuple", data.Array{data.String("a"), data.Map{"key1": data.Int(1)}, data.Array{data.Int(1), data.Int(2)}}},
			{"astuple", data.Array{data.String("a"), data.Map{"key1": data.Int(1)}, data.Array{data.Int(1), data.Int(2)}}},
		}

		for _, r := range returnTypes {
			r := r
			Convey(fmt.Sprint("When calling a function that returns ", r.typeName), func() {
				actual, err := mdl.Call(fmt.Sprintf("return_%s", r.typeName))
				So(err, ShouldBeNil)

				Convey("Then the function should return valid value", func() {
					So(actual, ShouldResemble, r.expected)
				})
			})
		}

		Reset(func() {
			mdl.Release()
		})
	})
}

func TestUnsupportedPyObject2Go(t *testing.T) {
	Convey("Given an initialized python py2go test module", t, func() {
		mainthread.AppendSysPath("")

		mdl, err := LoadModule("_test_py2go")
		So(err, ShouldBeNil)
		So(mdl, ShouldNotBeNil)

		Convey("When call python function which returns unsupported object", func() {
			ac, err := mdl.Call("return_object")

			Convey("Then converter should return an error", func() {
				So(ac, ShouldResemble, data.Null{})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "FailureTest")
			})
		})
	})
}
