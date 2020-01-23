package mainline

import (
	"encoding"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var SliceDelimiter = ","
var TimeFormat = time.RFC3339

// Parse the given strings into values of the given types.
// The length of both slices must be equal, with the type for the first string being the first type and so on.
func ValuesFromString(v []string, types []reflect.Type) ([]interface{}, error) {
	if len(v) < len(types) {
		return nil, fmt.Errorf("not enough parameters")
	}
	if len(v) > len(types) {
		return nil, fmt.Errorf("too many parameters")
	}

	vals := make([]interface{}, len(types))
	for i, pt := range types {
		val, err := ValueFromString(v[i], pt)
		if err != nil { // failed to parse as correct type, not a match
			return nil, fmt.Errorf("parameter %v could not be parsed as a %v", v[1], pt.String())
		}
		vals[i] = reflect.ValueOf(val).Elem().Interface()
	}
	return vals, nil
}

// ValueFromString attempts to parse the given string, into the given type.
// If the string is parsable and the type is supported, the resulting value is returned as an interface.
// Most types are supported with the exception of channels, functions.
// struct's must support either the json.Unmarshaler or encoding.TextUnmarshaler interfaces.
// Special cases for structs: URL and Time both supported
// The argument string is passed to these to unmarshal into the struct.
// slices/arrays are parsed as comma delimited items. Change the SliceDelimiter for something else.
// All supported types can be used as item types of the array.
// Base types float, int, bool string are supported.
// Maps is a work in progress ;-)
func ValueFromString(v string, t reflect.Type) (interface{}, error) {
	switch t.Kind() {
	case reflect.Interface, reflect.Ptr:
		return ValueFromString(v, t.Elem())

	case reflect.Struct:
		return structureFromString(v, t)

	case reflect.Slice:
		return sliceFromString(v, t)

	case reflect.Map:
		return mapFromString(v, t)

	case reflect.Float64, reflect.Float32:
		return floatFromString(v, t)

	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		return intFromString(v, t)

	case reflect.Bool:
		return boolFromString(v, t)

	case reflect.String:
		return stringFromString(v, t)

	default:
		return nil, fmt.Errorf("%s types are not supported as command line arguments", t.String())
	}
}

func structureFromString(s string, t reflect.Type) (interface{}, error) {
	pStr := reflect.New(t)

	if t == reflect.TypeOf(url.URL{}) {
		u, err := url.Parse(s)
		if err != nil {
			return nil, fmt.Errorf("%s could not be read as a %s  %v", s, t.String(), err)
		}
		return u, nil
	}

	if t == reflect.TypeOf(time.Time{}) {
		u, err := time.Parse(TimeFormat, s)
		if err != nil {
			return nil, fmt.Errorf("%s could not be read as a %s  %v", s, t.String(), err)
		}
		return u, nil
	}

	// If supports json, treat argument as json string
	if t.Implements(reflect.TypeOf((json.Unmarshaler)(nil))) {
		err := json.Unmarshal([]byte(s), pStr.Interface())
		if err != nil {
			return nil, err
		}
		return pStr.Interface(), nil
	}

	// If supports textUnmarshal, unmarshal argument into new object
	if t.Implements(reflect.TypeOf((*encoding.TextUnmarshaler)(nil))) {
		tu, ok := pStr.Interface().(encoding.TextUnmarshaler)
		if !ok {
			panic("Supposed supported interface didn't cast into that interface")
		}
		err := tu.UnmarshalText([]byte(s))
		if err != nil {
			return nil, err
		}
		return pStr.Interface(), nil
	}

	return nil, fmt.Errorf("failed to unmarshal argument %s into paramter %s as that parameter does not support a supported unmarshalling interface."+
		"Must support, json.Unmarshaler or encoding.TextUnmarshaler", s, t)
}

func sliceFromString(s string, t reflect.Type) (interface{}, error) {
	ss := strings.Split(s, SliceDelimiter)
	sv := reflect.MakeSlice(t.Elem(), 0, len(ss))
	for _, sa := range ss {
		sav, err := ValueFromString(sa, t.Elem())
		if err != nil {
			return nil, fmt.Errorf("%s could not be read as a %s", sa, t.Elem().String())
		}
		sv = reflect.Append(sv, reflect.ValueOf(sav).Elem())
	}
	return sv.Interface(), nil
}

func mapFromString(s string, t reflect.Type) (interface{}, error) {
	panic("Maps are not supported yet")
}

func floatFromString(s string, t reflect.Type) (interface{}, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, fmt.Errorf("%s could not be read as a %s", s, t.String())
	}
	iv := reflect.New(t)
	iv.Elem().SetFloat(f)
	return iv.Interface(), nil
}

func intFromString(s string, t reflect.Type) (interface{}, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%s could not be read as a %s", s, t.String())
	}
	iv := reflect.New(t)
	iv.Elem().SetInt(i)
	return iv.Interface(), nil
}

func boolFromString(s string, t reflect.Type) (interface{}, error) {
	b, err := strconv.ParseBool(s)
	if err != nil {
		b = true // Special case for bools, default to true, when present.
	}
	bv := reflect.New(t)
	bv.Elem().SetBool(b)
	return bv.Interface(), nil
}

func stringFromString(s string, t reflect.Type) (interface{}, error) {
	sv := reflect.New(t)
	sv.Elem().SetString(s)
	return sv.Interface(), nil
}
