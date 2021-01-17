// Copyright 2020 Rob Gilham
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package reflection

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
	case reflect.Interface:
		return ValueFromString(v, t.Elem())

	case reflect.Ptr:
		v, err := ValueFromString(v, t.Elem())
		if err != nil {
			return nil, err
		}
		p := reflect.New(t.Elem())
		p.Elem().Set(reflect.ValueOf(v))
		return p.Interface(), nil

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
	if s == "" {
		return pStr.Elem().Interface(), nil
	}

	if t == reflect.TypeOf(url.URL{}) {
		u, err := url.Parse(s)
		if err != nil {
			return nil, fmt.Errorf("%s could not be read as a %s  %v", s, t.String(), err)
		}
		return *u, nil
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
	sv := reflect.MakeSlice(t, 0, len(ss))
	for _, sa := range ss {
		sav, err := ValueFromString(sa, t.Elem())
		if err != nil {
			return nil, fmt.Errorf("%s could not be read as a %s", sa, t.Elem().String())
		}
		sv = reflect.Append(sv, reflect.ValueOf(sav).Elem())
	}
	return sv.Interface(), nil
}

// Map is parsed as json
func mapFromString(s string, t reflect.Type) (interface{}, error) {
	mp := reflect.New(t)
	if s != "" {
		err := json.Unmarshal([]byte(s), mp.Interface())
		if err != nil {
			return nil, err
		}
	} else {
		mp.Elem().Set(reflect.MakeMap(t))
	}
	return mp.Interface(), nil
}

func floatFromString(s string, t reflect.Type) (interface{}, error) {
	var f float64
	if s != "" {
		fl, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fmt.Errorf("%s could not be read as a %s", s, t.String())
		}
		f = fl
	}
	iv := reflect.New(t)
	iv.Elem().SetFloat(f)
	return iv.Interface(), nil
}

func intFromString(s string, t reflect.Type) (interface{}, error) {
	// Special cases
	if t == reflect.TypeOf(time.Duration(0)) {
		var d time.Duration
		if s != "" {
			du, err := time.ParseDuration(s)
			if err != nil {
				return nil, fmt.Errorf("%s could not be read as a %s  %v", s, t.String(), err)
			}
			d = du
		}
		return &d, nil
	}

	var i int
	if s != "" {
		ii, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%s could not be read as a %s", s, t.String())
		}
		i = int(ii)
	}
	return i, nil
}

func boolFromString(s string, t reflect.Type) (interface{}, error) {
	b := true // Special case for bools, default to true, when present.
	if s != "" {
		bb, err := strconv.ParseBool(s)
		if err != nil {
			return nil, fmt.Errorf("%s could not be read as a %s", s, t.String())
		}
		b = bb
	}
	bv := reflect.New(t)
	bv.Elem().SetBool(b)
	return bv.Interface(), nil
}

func stringFromString(s string, t reflect.Type) (interface{}, error) {
	sv := reflect.New(t)
	sv.Elem().SetString(s)
	return sv.Elem().Interface(), nil
}
