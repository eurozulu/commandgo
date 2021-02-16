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
package commandgo

import (
	"fmt"
	"github.com/eurozulu/commandgo/flags"
	"github.com/eurozulu/commandgo/functions"
	"os"
	"path"
	"strings"
)

var GlobalFlags = flags.NewFlags(true)

// Commands maps one or more 'command' strings to methods and/or functions on a mapped struct.
type Commands map[string]interface{}

// Run attempts to call the mapped method or function, using the first given argument as the key to the command map.
// If the given key is found, the remaining arguments are parsed into flags and parameters before the mapped method/func is called.
func (cmds Commands) Run(args ...string) error {
	// strip leading arg if it's program name
	if len(args) > 0 && args[0] == os.Args[0] {
		args = args[1:]
	}
	if err := GlobalFlags.Apply(args...); err != nil {
		return err
	}
	// adjust the arguments with any global flags removed
	args = GlobalFlags.Parameters()

	// use first arg as the command, if it exists. (Can be empty, is an empty mapping exists)
	var arg string
	if len(args) > 0 {
		arg = args[0]
		args = args[1:]
	}
	cmd, ok := cmds.findCommand(arg)
	if !ok {
		if arg == "" {
			return fmt.Errorf("no command given.  specify a command: %s <command>", path.Base(os.Args[0]))
		}
		return fmt.Errorf("'%s' is not a known command", arg)
	}
	i, ok := cmds[cmd]
	if !ok {
		return fmt.Errorf("CONFIG ERROR: command '%s' (%s) is not mapped", arg, cmd)
	}
	if i == nil {
		return fmt.Errorf("CONFIG ERROR: command '%s' (%s) is mapped to a nil value", arg, cmd)
	}

	if IsHelpCommand(i) {
		return CallHelpCommand(i, cmds, args...)
	}
	if functions.IsMethod(i) {
		return functions.CallMethod(i, args...)
	}
	if functions.IsFunc(i) {
		return functions.CallFunc(i, args...)
	}
	return fmt.Errorf("CONFIG ERROR: %v is an unknown type of function or method", i)
}

// findCommand looks through the map keys in non case sensative search
// returns the case sensative key if found or empty if not present
func (cmds Commands) findCommand(arg string) (string, bool) {
	for k := range cmds {
		if strings.EqualFold(k, arg) {
			return k, true
		}
	}
	return "", false
}
