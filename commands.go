package commandgo

import (
	"commandgo/functions"
	"commandgo/valuetyping"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Commands map[string]interface{}

func (c Commands) RunArgs() ([]interface{}, error) {
	return c.Run(os.Args[1:]...)
}

func (c Commands) Run(args ...string) ([]interface{}, error) {
	// first perform all assignments at this level of command, removing 'consumed' args from the command line.
	assKeys := c.assignments()
	cargs, err := c.applyAssignments(assKeys, args...)
	if err != nil {
		return nil, err
	}

	// Establish the command key.
	// If first carg known, use it as key and remaining cargs are parameters
	// if not known, use default "" key and use whole cargs as params
	var cmdKey string
	if len(cargs) > 0 {
		var ok bool
		cmdKey, ok = c.commandKey(cargs[0])
		if ok {
			cargs = cargs[1:]
		}
	}
	cmd, ok := c[cmdKey]
	if !ok {
		return nil, fmt.Errorf("unknown command '%s'", cmdKey)
	}

	// Check if command is a sub map
	cm, ok := cmd.(Commands)
	if ok {
		return cm.Run(cargs...)
	}
	return functions.CallFunc(cmd, cargs...)
}

func (c Commands) applyAssignments(keys []string, args ...string) ([]string, error) {
	var unused []string
	for i, arg := range args {
		ki := indexString(arg, keys)
		if ki < 0 {
			unused = append(unused, arg)
			continue
		}
		v := c[keys[ki]]

		var val string
		// Special case for bool assignments.  No value required, default to true
		if valuetyping.IsKind(v, reflect.Bool) {
			val = strconv.FormatBool(true)
			// check following arg parses as bool, otherwise ignore it
			if i+1 < len(args) {
				_, err := strconv.ParseBool(args[i+1])
				if err == nil {
					// following arg is a bool value, so use that
					i++
					val = args[i]
				}
			}
		} else {
			// not bool assignment, must have value
			if i+1 >= len(args) {
				return nil, fmt.Errorf("'%s' has no value", arg)
			}
			i++
			arg = args[i]
		}

		if err := valuetyping.SetValue(v, val); err != nil {
			return nil, err
		}
	}
	return unused, nil
}

func (c Commands) commandKey(arg string) (string, bool) {
	for k := range c {
		if strings.EqualFold(k, arg) {
			return k, true
		}
	}
	return "", false
}

func (c Commands) assignments() []string {
	var keys []string
	for k, v := range c {
		vo := reflect.TypeOf(v)
		if vo.Kind() != reflect.Ptr || vo.Kind() == reflect.Func {
			continue
		}
		keys = append(keys, k)
	}
	return keys
}

func (c Commands) functions() []string {
	var keys []string
	for k, v := range c {
		vo := reflect.TypeOf(v)
		if vo.Kind() != reflect.Func {
			continue
		}
		keys = append(keys, k)
	}
	return keys
}

func indexString(s string, in []string) int {
	for i, ins := range in {
		if strings.EqualFold(ins, s) {
			return i
		}
	}
	return -1
}
