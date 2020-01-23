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
