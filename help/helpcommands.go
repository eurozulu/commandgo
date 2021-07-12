package help

import (
	"strings"
)

const (
	HelpFlagShort = "-?"
	HelpFlagFull  = "--help"
)

// HelpLibrary are the globally available HelpGroups
var HelpLibrary []*HelpSubject

// HelpRequested is a flag to indicate the command is requesting help, rather than execution of the command.
// ShowHelp flags can be mapped to this point which, when true, will redirect the execution to the help system.
var HelpRequested bool

// ShowHelp is the main entry point for help.
// the given name may be a subject name, command or flag.
// returns the specific text for which ever is found matching the given name.
func ShowHelp(cmd string, args ...string) []interface{} {
	hs, hi := findSubject(cmd)
	if hs == nil {
		hs, hi = findSubject(args[0])
	}
	if hs == nil {
		// no matching help subject found, display root help
		hs, _ = findSubject("main")
	}

	var result []interface{}
	if hi != nil {
		result = append(result, hi.String())
	}
	if hs != nil {
		if hi != nil {
			result = append(result, hs.StringShort())
		} else {
			result = append(result, hs.String())
		}
	}
	return result
}

func findSubject(name string) (*HelpSubject, *HelpItem) {
	for _, hs := range HelpLibrary {
		if strings.EqualFold(hs.Name, name) {
			return hs, nil
		}

		hi := findItem(name, hs.HelpItems)
		if hi != nil {
			return hs, hi
		}
	}
	return nil, nil
}

func findItem(name string, items []*HelpItem) *HelpItem {
	for _, hi := range items {
		if hi.IsName(name) {
			return hi
		}
	}
	return nil
}
