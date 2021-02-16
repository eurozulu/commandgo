package mainline

import (
	"bytes"
	"fmt"
	"github.com/eurozulu/mainline/functions"
	"reflect"
	"runtime"
	"sort"
	"strings"
)

type HelpCommand struct {
	CommandMap Commands
}

// Help lists all the avilable commands and their aliases
func (ch HelpCommand) Help(_ ...string) error {
	// Invert command map so functions are the keys and stringslice of all commands mapped to it.
	iMap := map[string][]string{}
	for k, iv := range ch.CommandMap {
		if reflect.TypeOf(iv) == reflect.TypeOf(HelpCommand.Help) {
			continue
		}
		fName := runtime.FuncForPC(reflect.ValueOf(iv).Pointer()).Name()
		cmds := iMap[fName]
		iMap[fName] = append(cmds, k)
	}

	buf := bytes.NewBuffer(nil)
	for k, v := range iMap {
		// Sort command so longest is title command
		sort.Slice(v, func(i, j int) bool {
			return len(v[i]) > len(v[j])
		})

		_, _ = fmt.Fprintf(buf, "%s", v[0])
		if len(v) > 1 {
			_, _ = fmt.Fprintf(buf, " (%s)\t", strings.Join(v[1:], ", "))
		} else {
			_, _ = fmt.Fprintf(buf, "\t\t\t")
		}
		_, _ = fmt.Fprintln(buf, k)
	}
	bc := strings.Split(buf.String(), "\n")
	sort.Strings(bc)
	fmt.Println(strings.Join(bc, "\n"))
	return nil
}

func IsHelpCommand(i interface{}) bool {
	if !functions.IsMethod(i) {
		return false
	}
	return reflect.PtrTo(reflect.TypeOf(i).In(0)) == reflect.TypeOf((*HelpCommand)(nil))
}

func CallHelpCommand(i interface{}, cmds Commands, args ...string) error {
	if !IsHelpCommand(i) {
		return fmt.Errorf("interface is not HelpCommand")
	}
	hc := reflect.New(reflect.TypeOf(i).In(0)).Interface().(*HelpCommand)
	hc.CommandMap = cmds
	return hc.Help(args...)
}
