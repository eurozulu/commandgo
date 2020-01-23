package mainline

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
)

// commands is a slice of all the registered commands.
var commands map[string]Command

// MustAddCommand Adds a new command and panics if there is an error.
func MustAddCommand(name string, cmd interface{}) {
	if err := AddCommand(name, cmd); err != nil {
		panic(err)
	}
}

// AddCommand add a new command mapping to the given method or function.
// name must be a unique name.
// method must be the func type or a method on a struct.  Only Methods support flags.
func AddCommand(name string, fun interface{}) error {
	cmdName := strings.ToLower(name)
	_, ok := commands[cmdName]
	if ok {
		return fmt.Errorf("a command with the name %s already exists", name)
	}

	fv := reflect.ValueOf(fun)
	if fv.Kind() != reflect.Func {
		return fmt.Errorf("command %s must map to a func type, not a %s", name, fv.Kind().String())
	}

	fName := runtime.FuncForPC(fv.Pointer()).Name()
	fn := FuncCommand{
		name:      fName,
		function:  fv,
		signature: NewSignature(fv.Type(), false),
	}

	var cmd Command
	if isFuncMethod(fName, fv) {
		fn.signature = NewSignature(fv.Type(), true) // Skip the first param

		cmd = &MethodCommand{
			FuncCommand: fn,
			structType:  fv.Type().In(0),
		}
	} else {
		cmd = fn
	}

	if commands == nil {
		commands = make(map[string]Command)
	}
	commands[name] = cmd
	return nil
}

// RemoveCommand removes any command with the given name or does nothing if name doesn't exist.
func RemoveCommand(name string) {
	delete(commands, name)
}

// CommandNames gets a list of all the registered commands.
func CommandNames() []string {
	n := make([]string, len(commands))
	i := 0
	for v := range commands {
		n[i] = v
		i++
	}
	return n
}

// RunCommandLine executes one of the pre-defined commands based on the first argument following the executable name in os.Args.
// If the command returns an error it prints the error to standard err.
func RunCommandLine() {
	result, err := RunCommand(os.Args[1:]...)
	if nil != err {
		if _, err = fmt.Fprintln(os.Stderr, err); err != nil {
			panic(err)
		}
	}

	for _, v := range result {
		if _, err = os.Stdout.WriteString(ValueToString(v)); err != nil {
			panic(err)
		}
		if _, err = os.Stdout.WriteString("\n"); err != nil {
			panic(err)
		}
	}
}

// RunCommand maps the given arguments to the respective command and executes it.
// The first argument is used as the primary command mapping.
func RunCommand(args ...string) ([]interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("requires at least one argument")
	}
	cmd, err := findCommandByName(args[0])
	if err != nil {
		return nil, err
	}
	return cmd.Call(args)
}

func findCommandByName(name string) (Command, error) {
	var matched []string

	lName := strings.ToLower(name)
	for k, cmd := range commands {
		if strings.EqualFold(k, name) {
			return cmd, nil
		}
		if strings.HasPrefix(k, lName) {
			matched = append(matched, k)
		}
	}
	if len(matched) == 1 {
		return commands[matched[0]], nil
	}

	// matched multiple commands, list them as an error
	bf := bytes.NewBuffer(nil)
	for i, m := range matched {
		if i > 0 {
			bf.WriteString(" or")
		}
		bf.WriteRune(' ')
		bf.WriteString(m)
	}
	return nil, fmt.Errorf("%s is not a known command.  Did you mean %s", name, bf.String())
}

// isFuncMethod works out if the given func value is a Method or a regular func
func isFuncMethod(name string, fn reflect.Value) bool {
	ft := fn.Type()

	// Must have a parent struct as first arg
	if ft.NumIn() < 1 {
		return false
	}
	pn := parentName(name)
	pt := ft.In(0)
	if pt.String() != pn {
		return false
	}
	// Check they are the same func
	mName := baseName(name)
	m, ok := pt.MethodByName(mName)
	if !ok {
		return false
	}

	return m.Func.Pointer() == fn.Pointer()
}

func parentName(n string) string {
	parentName := strings.Split(n, ".")
	if len(parentName) < 1 {
		return ""
	}
	return strings.Join(parentName[:len(parentName)-1], ".")
}
func baseName(n string) string {
	bName := strings.Split(n, ".")
	if len(bName) < 1 {
		return ""
	}
	return bName[len(bName)-1]
}
