package commando

import (
	"fmt"
	"reflect"
)

type Command struct {
	// The name of this command
	Name string

	structType reflect.Type
	method     reflect.Method
}

// Run is the execution point for a command.
// The given arguments, less the first 'command' argument, are converted into flags and parameters
// and the matching method for those parameters is called on the under lying struct.
func (cmd Command) Run(args []string) error {
	// Parse flags first to remove them, and their data, from the command line.
	flgs, cmdLine, err := cmd.parseFlags(args[1:])
	if err != nil {
		return err
	}

	sig := cmd.Signature()
	params, err := cmd.parseParams(cmdLine, sig)
	if err != nil {
		return fmt.Errorf("'%s' failed as %v.  Requires arguments: %s", cmd.Name, err, sig.String())
	}

	// Create the new instance to be called
	pStr := reflect.New(cmd.structType)
	if err := cmd.setFields(flgs, pStr); err != nil {
		return err
	}

	// Call the matched method with the values set it matched with.
	mv := pStr.Elem().MethodByName(cmd.method.Name)
	rvs := mv.Call(params)

	// Check if any of the returned value types was an error
	errT := reflect.TypeOf(err)
	for _, rv := range rvs {
		if rv.Type().AssignableTo(errT) {
			return fmt.Errorf(rv.String())
		}
	}
	return nil
}

// Signatures builds a Signatures for the command method
func (cmd Command) Signature() *Signature {
	return NewSignature(cmd.method)
}

// parseParams parses the given argument list into values typed according to the given signature.
// The length of params and ParamTypes in the Signature MUST be equal as arguments are mapped, one to one
// according to their position.
func (cmd Command) parseParams(args []string, s *Signature) ([]reflect.Value, error) {
	if len(args) < len(s.ParamTypes) {
		return nil, fmt.Errorf("not enough parameters")
	}
	if len(args) > len(s.ParamTypes) {
		return nil, fmt.Errorf("too many parameters")
	}

	vals := make([]reflect.Value, len(s.ParamTypes))
	for i, pt := range s.ParamTypes {
		v, err := ValueFromString(args[i], pt)
		if err != nil { // failed to parse as correct type, not a match
			return nil, fmt.Errorf("parameter %d could not be parsed as a %s", args[1], pt.String())
		}
		vals[i] = reflect.ValueOf(v).Elem()
	}
	return vals, nil
}
