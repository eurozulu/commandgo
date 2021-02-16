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
package functions

import (
	"fmt"
	"github.com/eurozulu/commandgo/values"
	"reflect"
)

// ParseParameters parses the given argument slice of strings into a list of Values of the correct type
// for the given Signature
func ParseParameters(sig *Signature, args []string) ([]reflect.Value, error) {
	var vals []reflect.Value
	for i, pt := range sig.ParamTypes {
		var val interface{}
		var err error
		// if last param and variadic, wrap final args into a single array
		if sig.IsVariadic && i == len(sig.ParamTypes)-1 {
			if i < len(args) { // optional params provided
				vps, err := variadicParams(args[i:], pt.Elem())
				if err != nil {
					return nil, err
				}
				vals = append(vals, vps...)
			}
			continue

		} else if i < len(args) {
			val, err = values.ValueFromString(args[i], pt)
		} else {
			return nil, fmt.Errorf("missing argument %d, requires a %s value", i+1, pt.String())
		}
		if err != nil {
			return nil, fmt.Errorf("argument %d, %v", i, err)
		}
		vals = append(vals, reflect.ValueOf(val))
	}
	if !sig.IsVariadic && len(vals) < len(args) {
		return nil, fmt.Errorf("too many arguments.  %v expected, found %v", sig.String(), args)
	}
	return vals, nil
}

func variadicParams(args []string, t reflect.Type) ([]reflect.Value, error) {
	vals := make([]reflect.Value, len(args))
	for i, arg := range args {
		val, err := values.ValueFromString(arg, t)
		if err != nil { // failed to parse as correct type, not a match
			return nil, fmt.Errorf("parameter %v could not be parsed as a %v", arg, t.String())
		}
		vals[i] = reflect.ValueOf(val)
	}
	return vals, nil

}
