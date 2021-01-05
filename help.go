package mainline

import (
	"fmt"
	"sort"
	"strings"
)

// ShowCommands lists all the avilable commands and their aliases
func ShowCommands(cmds Commands, args ...string) {
	var c []string
	for k := range cmds {
		ks := strings.Split(k, ",")
		s := ks[0]
		if strings.HasPrefix(s, "-") {
			if len(ks) < 2 {
				return
			}
			s = ks[1]
		}
		if len(ks) > 1 {
			s = strings.Join([]string{s, fmt.Sprintf("\t\t(%s)", strings.Join(ks[1:], ", "))}, "")
		}
		c = append(c, s)
	}
	sort.Strings(c)
	for _, cmd := range c {
		fmt.Println(cmd)
	}
}
