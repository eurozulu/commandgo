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

package mainline

import (
	"fmt"
	"reflect"
)

// Command defines a single command mapped to a string name.
// The Command encompasses the underlying function the command maps to and its Name is of that function.
// The function may be a regular func or a Method function.  Only Method functions support flags.
type Command interface {
	// The name of this command
	Name() string

	// Call is the execution point for a command.
	// The given arguments, less the first 'command' argument, are parsed into typed values before being passed to the commands function.
	Call(args []string) ([]interface{}, error)
}

// FuncCommand is a regular function command.  It has no flags.
type FuncCommand struct {
	name      string
	function  reflect.Value
	signature Signature
}

func (cmd FuncCommand) Name() string {
	return baseName(cmd.name)
}

func (cmd FuncCommand) Call(args []string) ([]interface{}, error) {
	params, err := ValuesFromString(args[1:], cmd.signature.ParamTypes)
	if err != nil {
		return nil, fmt.Errorf("'%s' failed as %v.  Requires arguments: %s", cmd.Name(), err, cmd.signature.String())
	}

	vals := make([]reflect.Value, len(params))
	for i, p := range params {
		vals[i] = reflect.ValueOf(p)
	}
	rvs := cmd.function.Call(vals)
	var result []interface{}

	// Check if any of the returned value types was an error
	errT := reflect.TypeOf((*error)(nil)).Elem()
	for _, rst := range rvs {
		if rst.Type().AssignableTo(errT) {
			if !rst.IsNil() {
				return nil, fmt.Errorf("%v", rst.Interface())
			}
			continue
		}
		result = append(result, rst.Interface())
	}
	return result, nil
}
