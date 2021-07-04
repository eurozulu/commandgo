package help

import (
	"bytes"
	"sort"
	"strings"
)

const (
	HelpFlagShort = "?"
	HelpFlagFull  = "help"
)

// HelpSubjects are the globally available help subjects
var HelpSubjects []*HelpSubject

func Subject(name string) *HelpSubject {
	for _, hs := range HelpSubjects {
		if strings.EqualFold(hs.Name, name) {
			return hs
		}
	}
	return nil
}

// Help is the main entry point for help.
// the given name may be a group name, command or flag.
// returns the specific text for which ever is found matching the given name.
func Help(name string) string {
	hs := Subject(name)
	if hs != nil {
		return hs.Comment
	}

	// Search through the commands and flags
	for _, hs := range HelpSubjects {
		k, ok := findCommandString(name, hs.Commands)
		if ok {
			return hs.Commands[k]
		}
		k, ok = findCommandString(name, hs.Flags)
		if ok {
			return hs.Flags[k]
		}
	}
	return ""
}

// ListCommands lists all the Available Flags and commands, in their respective subject groups.
func ListCommands() []byte {
	out := bytes.NewBuffer(nil)
	for i, hs := range HelpSubjects {
		if i > 0 {
			out.WriteString("\n\n")
		}

		if hs.Comment != "" {
			out.WriteString(hs.Comment)
			out.WriteString("\n\n")
		}

		if len(hs.Flags) > 0 {
			out.WriteString("Flags:\n")
			keys := mapKeysOrdered(hs.Flags)
			for _, k := range keys {
				out.WriteRune('-')
				out.WriteString(k)
				out.WriteRune('\n')
				out.WriteString(hs.Flags[k])
				out.WriteString("\n\n")
			}
		}
		if len(hs.Commands) > 0 {
			out.WriteString("Commands:\n")
			keys := mapKeysOrdered(hs.Commands)
			for _, k := range keys {
				out.WriteString(k)
				out.WriteString("\t\t")
				cmd := strings.SplitN(hs.Commands[k], "\n", 2)
				out.WriteString(cmd[0])
				out.WriteString("\n")
			}
		}
	}
	return out.Bytes()
}

func findCommandString(arg string, m map[string]string) (string, bool) {
	for k := range m {
		if strings.EqualFold(k, arg) {
			return k, true
		}
	}
	return "", false
}

func mapKeysOrdered(m map[string]string) []string {
	s := make([]string, len(m))
	var i int
	for k := range m {
		s[i] = k
		i++
	}
	sort.Slice(s, func(i, j int) bool {
		return strings.Compare(s[i], s[j]) < 0
	})
	return s
}
