package mainline

import (
	"fmt"
	"reflect"
	"strings"
)

// commandMap contains the individual commands mapped to respective command object
type commandMap map[string]*command

// command is a container of the command object and the method to call
type command struct {
	method    *reflect.Method
	cmdObject interface{}
}

// newCommandMap builds the command map, dividing the comma delimited keys into single entires, pointing to the same command
func newCommandMap(cmds Commands) (commandMap, error) {
	km := commandMap{}
	for k, i := range cmds {
		if i == nil {
			return nil, fmt.Errorf("config error: command %s.  command struct is nil", k)
		}

		v := reflect.ValueOf(i)
		if v.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("config error: command %s.  is not a pointer to a struct", k)
		}
		if v.IsNil() {
			return nil, fmt.Errorf("config error: command %s.  is not a pointer to a struct", k)
		}
		if v.Elem().Kind() != reflect.Struct {
			return nil, fmt.Errorf("config error: command %s.  is not a pointer to a struct", k)
		}

		ks := strings.Split(k, ",")
		mn := strings.TrimLeft(ks[0], "-")
		m := findMethod(v.Elem().Type(), mn)
		if m == nil {
			return nil, fmt.Errorf("config error: command %s. is not known to Command structure %s", mn, v.Type().Name())
		}

		cmd := command{
			method:    m,
			cmdObject: i,
		}
		for _, nk := range ks {
			nk = strings.TrimSpace(nk)
			if strings.HasPrefix(nk, "-") {
				continue
			}
			if _, ok := km[nk]; ok {
				return nil, fmt.Errorf("config error: command '%s' is declared more than once", nk)
			}
			km[nk] = &cmd
		}
	}
	return km, nil
}

// findMethod finds the actual method name of the given string, regardles of case
func findMethod(t reflect.Type, name string) *reflect.Method {
	mCount := t.NumMethod()
	for i := 0; i < mCount; i++ {
		m := t.Method(i)
		if !strings.EqualFold(m.Name, name) {
			continue
		}
		return &m
	}
	return nil
}
