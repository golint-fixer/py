package py

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/py.v0/mainthread"
	"gopkg.in/sensorbee/sensorbee.v0/data"
)

func TestConvertGo2PyObject(t *testing.T) {
	Convey("Given an initialized python go2py test module", t, func() {
		mainthread.AppendSysPath("")

		mdl, err := LoadModule("_test_go2py")
		So(err, ShouldBeNil)
		So(mdl, ShouldNotBeNil)

		Convey("When set an object", func() {
			type argAndExpected struct {
				arg      data.Value
				expected string
			}
			values := map[string]argAndExpected{
				"string": argAndExpected{data.String("test"), "test"},
				"int":    argAndExpected{data.Int(9), "9"},
				"float":  argAndExpected{data.Float(0.9), "0.9"},
				"true":   argAndExpected{data.True, "True"},
				"false":  argAndExpected{data.False, "False"},
				"null":   argAndExpected{data.Null{}, "None"},
			}
			for k, v := range values {
				v := v
				msg := fmt.Sprintf("Then function should return string value: %v", k)
				Convey(msg, func() {
					actual, err := mdl.Call("go2py_tostr", v.arg)
					So(err, ShouldBeNil)
					So(actual, ShouldEqual, v.expected)
				})
			}
		})

		Convey("When set a time value", func() {
			now := time.Now().UTC()
			Convey("Then function should return time as string type", func() {
				actual, err := mdl.Call("go2py_tostr", data.Timestamp(now))
				So(err, ShouldBeNil)
				retStr, err := data.AsString(actual)
				So(err, ShouldBeNil)
				parsed, err := time.Parse("2006-01-02 15:04:05.999999999", retStr)
				So(err, ShouldBeNil)
				So(parsed, ShouldResemble, now.Truncate(time.Microsecond)) // Python's datetime has microseconds precision
			})
		})

		Convey("When set a byte array", func() {
			b := data.Blob([]byte("ABC"))
			Convey("Then function should return string", func() {
				actual, err := mdl.Call("go2py_toutf8", b)
				So(err, ShouldBeNil)
				So(actual, ShouldEqual, "ABC")
			})
		})

		Convey("When set a empty byte array", func() {
			b := data.Blob([]byte(""))
			Convey("Then function should return empty string", func() {
				actual, err := mdl.Call("go2py_toutf8", b)
				So(err, ShouldBeNil)
				So(actual, ShouldEqual, "")
			})
		})

		Convey("When set map in map and map in array", func() {
			arg := data.Map{
				"string": data.String("test"),
				"map": data.Map{
					"instr": data.String("test2"),
				},
				"array": data.Array{
					data.String("array-test"), data.Int(55),
				},
			}
			actual, err := mdl.Call("go2py_mapinmap", arg)
			Convey("Then function should return valid values", func() {
				So(err, ShouldBeNil)
				So(actual, ShouldEqual, "test_test2_array-test_55")
			})
		})

		Convey("When set array in array and map", func() {
			arg := data.Array{
				data.Array{
					data.String("test"), data.Int(55),
				},
				data.Map{
					"map": data.String("inmap"),
				},
			}
			actual, err := mdl.Call("go2py_arrayinmap", arg)
			Convey("Then function should return valid values", func() {
				So(err, ShouldBeNil)
				So(actual, ShouldEqual, "test_55_inmap")
			})
		})

		Reset(func() {
			mdl.Release()
		})
	})
}
