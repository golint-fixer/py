package py

/*
#include "Python.h"

PyObject* getPyNone() {
  return Py_BuildValue("");
}
*/
import "C"
import (
	"fmt"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"unsafe"
)

func getNewPyDic(m map[string]interface{}) Object {
	return Object{}
}

func newPyObj(v data.Value) (Object, error) {
	var pyobj *C.PyObject
	var err error
	switch v.Type() {
	case data.TypeBool:
		b, _ := data.ToInt(v)
		pyobj = C.PyBool_FromLong(C.long(b))
	case data.TypeInt:
		i, _ := data.AsInt(v)
		pyobj = C.PyLong_FromLong(C.long(i))
	case data.TypeFloat:
		f, _ := data.AsFloat(v)
		pyobj = C.PyFloat_FromDouble(C.double(f))
	case data.TypeString:
		s, _ := data.AsString(v)
		pyobj = newPyString(s)
	case data.TypeBlob:
		b, _ := data.AsBlob(v)
		cb := (*C.char)(unsafe.Pointer(&b[0]))
		pyobj = C.PyBytes_FromStringAndSize(cb, C.Py_ssize_t(len(b)))
	case data.TypeTimestamp:
		t, _ := data.AsTimestamp(v)
		pyobj = getPyDateTime(t)
	case data.TypeArray:
		innerArray, _ := data.AsArray(v)
		pyobj, err = newPyArray(innerArray)
	case data.TypeMap:
		innerMap, _ := data.AsMap(v)
		pyobj, err = newPyMap(innerMap)
	case data.TypeNull:
		pyobj = C.getPyNone()
	default:
		err = fmt.Errorf("unsupported type in sensorbee/py: %s", v.Type())
	}

	if pyobj == nil && err == nil {
		return Object{}, getPyErr()
	}
	return Object{p: pyobj}, err
}

func newPyString(s string) *C.PyObject {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return C.PyUnicode_FromString(cs)
}

func newPyArray(a data.Array) (*C.PyObject, error) {
	pylist := C.PyList_New(C.Py_ssize_t(len(a)))
	if pylist == nil {
		return nil, getPyErr()
	}
	for i, v := range a {
		value, err := newPyObj(v)
		if err != nil {
			return nil, err
		}
		// PyList object takes over the value's reference, and not need to
		// decrease reference counter.
		C.PyList_SetItem(pylist, C.Py_ssize_t(i), value.p)
	}
	return pylist, nil
}

func newPyMap(m data.Map) (*C.PyObject, error) {
	pydict := C.PyDict_New()
	if pydict == nil {
		return nil, getPyErr()
	}
	for k, v := range m {
		err := func() error {
			ck := C.CString(k)
			defer C.free(unsafe.Pointer(ck))
			value, err := newPyObj(v)
			if err != nil {
				return err
			}
			defer value.decRef()
			C.PyDict_SetItemString(pydict, ck, value.p)
			return nil
		}()
		if err != nil {
			return nil, err
		}
	}
	return pydict, nil
}
