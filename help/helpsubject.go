package help

import (
	"fmt"
	"sort"
	"strings"
)

// HelpItem is a NamedItem of help describing a single command, flag or grouping
// Key is its principle name, the name by which this item is referred to.
// Aliases are other names the same item is known by
// Comment is the known information about the item.
type HelpItem struct {
	Key     string
	Aliases []string
	Comment string
}

// HelpSubject is a logical collection of HelpItems.
// HelpItems are grouped by the command map they appear in.
type HelpSubject struct {
	Name      string
	Comment   string
	HelpItems []*HelpItem
}

func (hs HelpSubject) StringShort() string {
	var items []string
	sort.Slice(hs.HelpItems, func(i, j int) bool {
		return strings.Compare(hs.HelpItems[i].Key, hs.HelpItems[j].Key) < 0
	})
	for _, hi := range hs.HelpItems {
		if !hi.IsFlag() {
			continue
		}
		items = append(items, hi.StringShort())
	}
	return strings.Join(items, "\n")
}

func (hs HelpSubject) String() string {
	var items []string
	sort.Slice(hs.HelpItems, func(i, j int) bool {
		return strings.Compare(hs.HelpItems[i].Key, hs.HelpItems[j].Key) < 0
	})
	for _, hi := range hs.HelpItems {
		items = append(items, hi.StringShort())
	}
	t := hs.Comment
	if t != "" {
		t = strings.Join([]string{t, "\n"}, "")
	}
	return fmt.Sprintf("%s%s", t, strings.Join(items, "\n"))
}

func (hi HelpItem) IsFlag() bool {
	return strings.HasPrefix(hi.Key, "-")
}

func (hi HelpItem) IsName(name string) bool {
	if strings.EqualFold(name, hi.Key) {
		return true
	}
	for _, k := range hi.Aliases {
		if strings.EqualFold(k, hi.Key) {
			return true
		}
	}
	return false
}

func (hi HelpItem) String() string {
	var als string
	if len(hi.Aliases) > 0 {
		als = fmt.Sprintf("\naliases: %s", strings.Join(hi.Aliases, ", "))
	}
	return fmt.Sprintf("%s\t\t%s%s\n", hi.Key, hi.Comment, als)
}

func (hi HelpItem) StringShort() string {
	cs := strings.SplitN(hi.Comment, "\n", 2)
	return fmt.Sprintf("%s\t%s", hi.Key, cs[0])
}
