// Copyright 2020 Rob Gilham
//
// Licensed under the Apache License, Version newtype.0 (the "License");
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

package values

import (
	"encoding"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// SliceDelimiter determines how argument lists (contained within a single argument) are split.
var SliceDelimiter = ","

var textUnmarshalerInterface = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
var binaryUnmarshalerInterface = reflect.TypeOf((*encoding.BinaryUnmarshaler)(nil)).Elem()

// ValueFromString attempts to parse the given string, into the given type.
// If the string is parsable and the type is supported, the resulting value is returned as an interface.
// Most types are supported with the exception of channels, functions.
// All the Base types float32/64, int8,16,...64, bool string are supported.
// strings can be parsed into more complext structures:
// struct's may support either the encoding.TextUnmarshaler or encoding.BinaryUnmarshaler interfaces, otherwise attempts to parge argument as json
// The argument string is passed to these to unmarshal into the struct.
// slices/arrays are parsed as comma delimited items. Change the SliceDelimiter for something else.
// All supported types can be used as item types of the array.
// Maps are parsed as json structures. e.g. -mapflag '{"mykey": "myvalue", "isIt": true}'
// see CustomType to add additional types as valid parameter types.
func ValueFromString(v string, t reflect.Type) (interface{}, error) {
	// Check custom types first
	if IsCustomType(t) {
		return CustomValueFromString(v, t)
	}

	switch t.Kind() {
	case reflect.Ptr:
		// If type is a ptr, get value of element and wrap it in a new pointer
		vp, err := ValueFromString(v, t.Elem())
		if err != nil {
			return nil, err
		}
		p := reflect.New(t.Elem())
		p.Elem().Set(reflect.ValueOf(vp))
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

	case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint:
		return uintFromString(v, t)

	case reflect.Bool:
		return boolFromString(v, t)

	case reflect.String:
		return stringFromString(v, t)

	default:
		return nil, fmt.Errorf("%s types are not supported as command line arguments", t.String())
	}
}

func IsKind(i interface{}, k reflect.Kind) bool {
	if i == nil {
		return false
	}
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Kind() == k
}

func GetValue(r interface{}) interface{} {
	t := reflect.TypeOf(r)
	if t.Kind() == reflect.Ptr {
		return GetValue(reflect.ValueOf(r).Elem().Interface())
	}
	return reflect.ValueOf(r).Interface()
}

// Sets the given receiver with the given value.
// Assigns the value or a pointer to it, depending on the reciever type
func SetValue(r interface{}, val string) error {
	iVal, err := ValueFromString(val, reflect.TypeOf(r))
	if err != nil {
		return err
	}

	recv := reflect.ValueOf(r)
	if recv.Type().Kind() == reflect.Ptr {
		recv = recv.Elem()
	}

	v := reflect.ValueOf(iVal)
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}
	recv.Set(v)
	return nil
}

func structureFromString(s string, t reflect.Type) (interface{}, error) {
	pStr := reflect.New(t)
	if s == "" {
		return pStr.Elem().Interface(), nil
	}
	// If supports textUnmarshal, unmarshal argument into new object
	if t.Implements(textUnmarshalerInterface) {
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

	// If supports BinaryUnmarshaler, unmarshal argument into new object
	if t.Implements(binaryUnmarshalerInterface) {
		tu, ok := pStr.Interface().(encoding.BinaryUnmarshaler)
		if !ok {
			panic("Supposed supported interface didn't cast into that interface")
		}
		err := tu.UnmarshalBinary([]byte(s))
		if err != nil {
			return nil, err
		}
		return pStr.Interface(), nil
	}
	// try to parse as json
	err := json.Unmarshal([]byte(s), pStr.Interface())
	if err != nil {
		return nil, err
	}
	return pStr.Elem().Interface(), nil
}

func sliceFromString(s string, t reflect.Type) (interface{}, error) {
	ss := strings.Split(s, SliceDelimiter)
	sv := reflect.MakeSlice(t, 0, len(ss))
	for _, sa := range ss {
		sel, err := ValueFromString(strings.TrimSpace(sa), t.Elem())
		if err != nil {
			return nil, fmt.Errorf("%s could not be read as a %s", sa, t.Elem().String())
		}
		ev := reflect.ValueOf(sel)
		if ev.Kind() == reflect.Ptr {
			ev = ev.Elem()
		}
		sv = reflect.Append(sv, ev)
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
	return mp.Elem().Interface(), nil
}

func floatFromString(s string, t reflect.Type) (interface{}, error) {
	v := reflect.New(t)
	if s != "" {
		fl, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fmt.Errorf("%s could not be read as a %s", s, t.String())
		}
		if t.Kind() == reflect.Float32 && fl > math.MaxFloat32 {
			return nil, fmt.Errorf("argument %s could not be parsed as a %s", s, t.Name())
		}
		v.Elem().SetFloat(fl)
	}
	return v.Elem().Interface(), nil
}

func intFromString(s string, t reflect.Type) (interface{}, error) {
	// Special case for time.Duration
	// TODO: Can this be a custom type?
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

	v := reflect.New(t)
	if s != "" {
		ii, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%s could not be read as a %s", s, t.String())
		}
		v.Elem().SetInt(ii)
		// Ensure values match, in case of overflow.
		if v.Elem().Int() != ii {
			return nil, fmt.Errorf("argument %s could not be parsed as a %s", s, t.Name())
		}

	}
	return v.Elem().Interface(), nil
}

func uintFromString(s string, t reflect.Type) (interface{}, error) {
	v := reflect.New(t)
	if s != "" {
		ii, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%s could not be read as a %s", s, t.String())
		}
		v.Elem().SetUint(ii)
		// Ensure values match, in case of overflow.
		if v.Elem().Uint() != ii {
			return nil, fmt.Errorf("argument %s could not be parsed as a %s", s, t.Name())
		}
	}
	return v.Elem().Interface(), nil
}

func boolFromString(s string, t reflect.Type) (interface{}, error) {
	b := true // Special case for bools, default to true
	if s != "" {
		bb, err := strconv.ParseBool(s)
		if err != nil {
			return nil, fmt.Errorf("%s could not be read as a %s", s, t.String())
		}
		b = bb
	}
	return b, nil
}

func stringFromString(s string, t reflect.Type) (interface{}, error) {
	// remember its a string Kind but maybe not a string Type so we still create a new instance of the type and set the arg string
	sv := reflect.New(t)
	sv.Elem().SetString(s)
	return sv.Elem().Interface(), nil
}
