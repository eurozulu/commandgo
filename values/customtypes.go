// CustomTypes allow the type assertion to be expanded to support any type required.
// The custom typeing maps a Type to a specific value parsing function.
// The type defines which parameter types or variable can be used as mapping points, to receive a string argument.
// Any parameter which the type is assignable to can be specified as a mapped point (variable or function parameter)
// e.g. By specifying the *os.File as a type, a function with a signature:  myfunc(in io.ReadCloser) error
// can be mapped, and will recieve the new *File instance when called.
// The Type is mapped to a ArgValue function `func(s string) (interface{}, error)`
// This receives the argument from the command line and should 'parse' it into the required type.
// e.g. for the *File example, (Which is built in already) the argument is treated as the file path and
// a call to os.Open is made, returning the resulting file.

package values

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"time"
)

// ArgValue parses the single string argument into a specific type.
type ArgValue func(s string, t reflect.Type) (interface{}, error)

// Timeformat for which time types are parsed
var TimeFormat = time.RFC3339

var customTypes = map[reflect.Type]ArgValue{}

// NewCustomType adds the given type as a new, valid parameter type, which can be parsed from string by the given ArgValue function.
// to remove a mapping, add the type with a nil value.
func NewCustomType(t reflect.Type, pfunc ArgValue) {
	if pfunc == nil {
		if _, ok := customTypes[t]; ok {
			delete(customTypes, t)
		}
		return
	}
	customTypes[t] = pfunc
}

// IsCustomType checks if the given type will be supported
func IsCustomType(t reflect.Type) bool {
	return customType(t) != nil
}

func CustomValueFromString(arg string, t reflect.Type) (interface{}, error) {
	cv := customType(t)
	if cv == nil {
		return nil, fmt.Errorf("unsupported custome type")
	}
	return cv(arg, t)
}

func customType(t reflect.Type) ArgValue {
	for k, v := range customTypes {
		if k.AssignableTo(t) {
			return v
		}
	}
	return nil
}

// init registers the "out of the box" custom types supported
func init() {
	NewCustomType(reflect.TypeOf(&os.File{}), customTypeFile)
	NewCustomType(reflect.TypeOf(&url.URL{}), customTypeURL)
	NewCustomType(reflect.TypeOf(time.Time{}), customTypeTime)
}

func customTypeFile(s string, t reflect.Type) (interface{}, error) {
	if s == "-" {
		return os.Stdin, nil
	}
	return os.Open(s)
}

func customTypeURL(s string, t reflect.Type) (interface{}, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("%s could not be read as a url  %v", s, err)
	}
	return u, nil
}

func customTypeTime(s string, t reflect.Type) (interface{}, error) {
	u, err := time.Parse(TimeFormat, s)
	if err != nil {
		return nil, err
	}
	return u, nil
}
