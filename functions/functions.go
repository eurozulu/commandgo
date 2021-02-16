package functions

import (
	"fmt"
	"github.com/eurozulu/mainline/flags"
	"reflect"
	"runtime"
	"strings"
)

// Checks if given interface is a func.
// will be true for both global functions and methods.
func IsFunc(i interface{}) bool {
	return reflect.TypeOf(i).Kind() == reflect.Func
}

// CallFunc calls the given function interface using the given arguments.
// interface must be a fucntion (IsFunc returns true).
// function is called as a global function, assuming all parameters are inputs.
// If called with a method, will assume the reciever sturcture is a parameter.
func CallFunc(i interface{}, args ...string) error {
	if !IsFunc(i) {
		return fmt.Errorf("Not a function")
	}

	// Check for unknown flags
	fgs := flags.NewFlags(false)
	if err := fgs.Apply(args...); err != nil {
		return err
	}

	// parse args into parameters
	sig, err := NewSignature(i)
	if err != nil {
		return err
	}
	inVals, err := ParseParameters(sig, args)
	if err != nil {
		return err
	}
	outVals := reflect.ValueOf(i).Call(inVals)

	// check if an error returned
	// TODO: review how to handle non error return values, if at all.
	errInterface := reflect.TypeOf((*error)(nil)).Elem()
	for _, ov := range outVals {
		if (ov.Kind() != reflect.Interface) || ov.IsNil() {
			continue
		}
		if ov.Type().Implements(errInterface) {
			return ov.Interface().(error)
		}
	}
	return nil
}

// Checks if the given interface is a method.
// Must be a function (IsFunc returns true) AND have a first parameter of a struct reciever.
// The reciever must have a method with the func name
func IsMethod(i interface{}) bool {
	if !IsFunc(i) {
		return false
	}
	// func must have first param as struct
	vt := reflect.TypeOf(i)
	if vt.NumIn() < 1 {
		return false
	}
	p1 := vt.In(0)
	if p1.Kind() != reflect.Struct {
		return false
	}
	// ensure that struct and func name match. (Not just a random struct as a parameter)
	if _, ok := p1.MethodByName(FuncName(i, false)); !ok {
		return false
	}
	return true
}

func CallMethod(i interface{}, args ...string) error {
	if !IsMethod(i) {
		return fmt.Errorf("Not a method!")
	}

	// new instance of struct and get ref to method
	ns := reflect.New(reflect.TypeOf(i).In(0))
	md := ns.MethodByName(FuncName(i, false))

	flgs, err := flags.NewStructFlags(ns.Type())
	if err != nil {
		return err
	}
	if err := flgs.Apply(args...); err != nil {
		return err
	}

	// parse args into parameters
	sig, err := NewSignature(i)
	if err != nil {
		return err
	}
	inVals, err := ParseParameters(sig, flgs.Parameters())
	if err != nil {
		return err
	}

	outVals := md.Call(inVals)

	// check if an error returned
	errInterface := reflect.TypeOf((*error)(nil)).Elem()
	for _, ov := range outVals {
		if (ov.Kind() != reflect.Interface) || ov.IsNil() {
			continue
		}
		if ov.Type().Implements(errInterface) {
			return ov.Interface().(error)
		}
	}
	return nil
}

// Get the function name if the given interface is a func.
// If not a func or is nil, , returns empty string
// withPackage flag, when true privades a dot delmited <package>.<name>
// when false, just the function name itself (string after last ".")
func FuncName(i interface{}, withPackage bool) string {
	if !IsFunc(i) {
		return ""
	}
	v := reflect.ValueOf(i)
	if v.IsNil() {
		return ""
	}
	fn := runtime.FuncForPC(v.Pointer()).Name()
	if withPackage {
		return fn
	}
	fns := strings.Split(fn, ".")
	if len(fns) == 0 {
		return fn
	}
	return fns[len(fns)-1]
}
