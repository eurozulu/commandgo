package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/eurozulu/mainline"
	"io"
	"os/exec"
	"strings"
)

// MappedCommands is a map of the functions which have been mapped to a single command
type MappedCommands map[string]*CommandDetails

// CommandDetails represents a single command mapping to its function
type CommandDetails struct {
	FullName   string `json:"packagename"`
	Parameters string `json:"parameters"`
	Comment    string `json:"comment"`
}

// ParseCommands parses the go doc output to extract the relevant comments from the function which have been mapped.
func ParseCommands(pkg string, out io.Writer) error {
	all, err := RunDoc("-src", "-all", "-u")
	if err != nil {
		return err
	}

	mcs := findMappedCommands(all)
	for _, mc := range mcs {
		gd, err := RunDoc(mc.FullName)
		if err != nil {
			panic(err)
		}
		lns := strings.SplitN(gd, "\n", 2)
		mc.Parameters = readParameters(lns[0])
		mc.Comment = strings.TrimSpace(strings.Join(lns[1:], "\n"))
	}

	return json.NewEncoder(out).Encode(mcs)
}

func findMappedCommands(cmd string) MappedCommands {
	scn := bufio.NewScanner(bytes.NewBuffer([]byte(cmd)))
	cmds := []string{
		mainline.FuncName(mainline.MustAddCommand),
		mainline.FuncName(mainline.AddCommand),
	}

	mappings := make(MappedCommands)

	for scn.Scan() {
		ln := scn.Text()
		if ln == "" || strings.HasPrefix(ln, "//") {
			continue
		}
		for _, c := range cmds {
			if strings.Contains(ln, c) {
				pm := readParameters(ln)
				s := strings.Split(pm, ",")
				key := strings.Trim(s[0], "'\"")
				mappings[key] = &CommandDetails{
					FullName: strings.TrimSpace(s[1]),
				}
				break
			}
		}
	}
	return mappings
}

// Reads the contents of any string enclosed in parameters
func readParameters(ln string) string {
	i := strings.Index(ln, "(")
	if i < 0 || i >= len(ln) {
		return ""
	}
	ie := strings.Index(ln[i:], ")")
	if ie < 0 {
		return ""
	}
	return ln[i+1 : i+ie]
}

func RunDoc(args ...string) (string, error) {
	buf := bytes.NewBuffer(nil)
	args = append([]string{"doc"}, args...)
	c := exec.Command("go", args...)
	c.Stdout = buf

	if err := c.Run(); err != nil {
		return "", err
	}
	return buf.String(), nil
}
