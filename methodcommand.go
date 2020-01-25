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
	"strconv"
	"strings"
)

const tagKey = "flag"

type MethodCommand struct {
	FuncCommand
	structType reflect.Type
}

// Flag is a named argument, preceeded with a dash or dounble dash. '-', '--'
// All flags have a name, which is the remainder of the string following the dash(es).
// With the exception of boolean flags, all flags must have a value, the argument following the name, separated with a space.
// Boolean flags may have a value, if the following argument is (present and) parsable via #strconv.ParseBool, otherwise they are assumed to be true
// and any following argument is treated as a parameter.
// A Flag contains a Value as an interface, which will be a value of the type for that Field.  (ie. an int Field has an int value)
// The FieldIndex is the index of the Field within its parent struct.
type Flag struct {
	Name       string
	Value      interface{}
	FieldIndex int
}

func (cmd MethodCommand) Call(args []string) ([]interface{}, error) {
	// Parse flags first to remove them, and their data, from the command line.
	flgs, cmdLine, err := cmd.parseFlags(args[1:])
	if err != nil {
		return nil, err
	}

	// Create the new instance of the parent struct
	pStr := reflect.New(cmd.structType)
	if err := cmd.setFields(flgs, pStr); err != nil {
		return nil, err
	}

	// Call the matched method with the values set it matched with.
	cmd.FuncCommand.function = pStr.Elem().MethodByName(cmd.Name())

	cmdLine = append([]string{args[0]}, cmdLine...) // Insert the command again as function will strip it off
	return cmd.FuncCommand.Call(cmdLine)
}

// setFields sets all the given flags on the given (pointer) struct's fields.
func (cmd MethodCommand) setFields(flds []Flag, pStr reflect.Value) error {
	for _, f := range flds {
		fld := pStr.Elem().Field(f.FieldIndex)
		if !fld.IsValid() || !fld.CanSet() {
			return fmt.Errorf("flag %s is not an addressable field", f.Name)
		}

		// Assign field value.  Check if receiver is expecting a pointer or not.
		var v reflect.Value
		if fld.Type().Kind() == reflect.Ptr {
			v = reflect.ValueOf(f.Value)
		} else {
			v = reflect.ValueOf(f.Value).Elem()
		}
		fld.Set(v)
	}
	return nil
}

// parseFlags parses the given arguments for flags
// returns the flags as a slice and the remaining, unnamed, arguments as a slice of strings
// Flag values are converted into their relevant type for the corrisponding field they are mapped to.
func (cmd MethodCommand) parseFlags(args []string) ([]Flag, []string, error) {
	var flgs []Flag
	var params []string

	for i := 0; i < len(args); i++ {
		if !strings.HasPrefix(args[i], "-") {
			params = append(params, args[i])
			continue
		}
		name := strings.TrimLeft(args[i], "-")
		inx := cmd.targetFieldIndex(name)
		if inx < 0 {
			return nil, nil, fmt.Errorf("--%s is not a known flag", name)
		}

		fld := cmd.structType.Field(inx)
		var val interface{}
		// special case for bools as they dont need the next argument
		if fld.Type.Kind() == reflect.Bool {
			b := true
			if i+1 < len(args) { // if there is a following args, check if it can be parsed as bool
				bt, err := strconv.ParseBool(args[i+1])
				if err == nil {
					b = bt
					i++ // move past the 'consumed' argument
				}
			}
			val = &b

		} else { // All other, non bool types
			if i+1 >= len(args) {
				return nil, nil, fmt.Errorf("missing value for %s flag", name)
			}
			var err error
			val, err = ValueFromString(args[i+1], fld.Type)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse flag '%s' %v", name, err)
			}
			i++ // skip past the value arg
		}

		flgs = append(flgs, Flag{
			Name:       name,
			Value:      val,
			FieldIndex: inx,
		})
	}
	return flgs, params, nil
}

// targetFieldIndex finds the index of the field in the target struct, of the given name.
// name is checked, first against any "flag" tag.  If not present, checks for name match.
func (cmd MethodCommand) targetFieldIndex(name string) int {
	for i := 0; i < cmd.structType.NumField(); i++ {
		fd := cmd.structType.Field(i)
		tag, ok := fd.Tag.Lookup(tagKey)
		if !ok { // no tag, use field name
			tag = fd.Name
		}
		if tag == "-" { // ignore those taged with dash
			continue
		}
		names := strings.Split(tag, ",")
		for _, n := range names {
			if strings.EqualFold(name, strings.TrimSpace(n)) {
				return i
			}
		}
	}
	return -1
}
