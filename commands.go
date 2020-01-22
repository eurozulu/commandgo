package mainline

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// commands is a slice of all the registered commands.
var commands map[string]*Command

// MustAddCommand Adds a new command and panics if there is an error.
func MustAddCommand(name string, cmd interface{}) {
	if err := AddCommand(name, cmd); err != nil {
		panic(err)
	}
}

// AddCommand add a new command mapping to the given method.
// name must be a unique name.
// method must be the func type of a method (NOT regular function)
func AddCommand(name string, method interface{}) error {
	if commands == nil {
		commands = make(map[string]*Command)
	}

	name = strings.ToLower(name)
	_, ok := commands[name]
	if ok {
		return fmt.Errorf("a command with the name %s already exists", name)
	}

	st, mt, err := methodFromFunc(method)
	if err != nil {
		return err
	}

	commands[name] = &Command{
		Name:       name,
		method:     *mt,
		structType: st,
	}
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

func RunCommandLine() error {
	return RunCommand(os.Args[1:]...)
}

// RunCommand maps the given arguments to the respective command and executes it.
// The first argument is used as the primary command mapping.
func RunCommand(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("requires at least one argument")
	}
	cmd, err := findCommandByName(args[0])
	if err != nil {
		return err
	}
	return cmd.Run(args)
}

func findCommandByName(name string) (*Command, error) {
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

// methodFromFunc locates the parent structure and the method whic the given function points to.
func methodFromFunc(f interface{}) (reflect.Type, *reflect.Method, error) {
	ft := reflect.TypeOf(f)
	if ft.Kind() != reflect.Func {
		return nil, nil, fmt.Errorf("command function must be a func type, not a %s", ft.Kind().String())
	}
	if ft.NumIn() < 1 {
		return nil, nil, fmt.Errorf("command function must be a method on a struct, not a stand alone function")
	}
	st := ft.In(0)

	fv := reflect.ValueOf(f)
	fa := fmt.Sprintf("%v", fv)

	mc := st.NumMethod()
	for i := 0; i < mc; i++ {
		m := st.Method(i)
		id := fmt.Sprintf("%v", m.Func.Interface())
		if id == fa {
			return st, &m, nil
		}
	}
	return st, nil, fmt.Errorf("Failed to find func in parent structure.")
}
