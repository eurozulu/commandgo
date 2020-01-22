package commando

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// Commands is a collection of command structs to execute a command line argument against.
type Commands struct {
	// commands is a slice of all the registered commands.
	commands map[string]*Command
}

// MustAddCommand Adds a new command and panics if there is an error.
func (c *Commands) MustAddCommand(name string, cmd interface{}) {
	if err := c.AddCommand(name, cmd); err != nil {
		panic(err)
	}
}

func (c *Commands) AddCommand(name string, method interface{}) error {
	name = strings.ToLower(name)
	_, ok := c.commands[name]
	if ok {
		return fmt.Errorf("a command with the name %s already exists", name)
	}

	st, mt, err := c.methodFromFunc(method)
	if err != nil {
		return err
	}

	if c.commands == nil {
		c.commands = make(map[string] *Command)
	}

	c.commands[name] = &Command{
		Name:       name,
		method:     *mt,
		structType: st,
	}
	return nil
}

// RemoveCommand removes any command with the given name or does nothing if name doesn't exist.
func (c Commands) RemoveCommand(name string) {
	delete(c.commands, name)
}

// CommandNames gets a list of all the registered commands.
func (c Commands) CommandNames() []string {
	n := make([]string, len(c.commands))
	i := 0
	for v := range c.commands {
		n[i] = v
		i++
	}
	return n
}

func (c Commands) RunCommandLine() error {
	return c.RunCommand(os.Args[1:]...)
}

// RunCommand maps the given arguments to the respective command and executes it.
// The first argument is used as the primary command mapping.
func (c Commands) RunCommand(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("requires at least one argument")
	}
	cmd, err := c.findCommandByName(args[0])
	if err != nil {
		return err
	}
	return cmd.Run(args)
}

func (c Commands) findCommandByName(name string) (*Command, error) {
	var matched []string

	lName := strings.ToLower(name)
	for k, cmd := range c.commands {
		if strings.EqualFold(k, name) {
			return cmd, nil
		}
		if strings.HasPrefix(k, lName) {
			matched = append(matched, k)
		}
	}
	if len(matched) == 1 {
		return c.commands[matched[0]], nil
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
func (c Commands) methodFromFunc(f interface{}) (reflect.Type, *reflect.Method, error) {
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
