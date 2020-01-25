package helper

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/eurozulu/mainline"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

const helpdir = "_help"
const helpJSON = "help.json"
const helpGo = "help.go"

// MappedCommands is a map of the functions which have been mapped to a command
type MappedCommands map[string]*CommandDetails

// CommandDetails represents a single command mapping to its function
type CommandDetails struct {
	FullName   string `json:"packagename"`
	Parameters string `json:"parameters"`
	HelpText   string `json:"helptext"`
}

func GenHelp() error {
	p, _ := os.Getwd()
	hd := path.Join(p, helpdir)

	// Clean out existing
	if err := os.RemoveAll(hd); err != nil {
		return err
	}
	if err := os.MkdirAll(hd, 0755); err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	if err := readJsonHelp(buf); err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(hd, helpJSON), buf.Bytes(), 0644)
}

func CleanJson() error {
	p, _ := os.Getwd()
	hd := path.Join(p, helpdir, helpJSON)
	return os.RemoveAll(hd)
}

func readJsonHelp(out io.Writer) error {
	all, err := RunDoc("-src", "-all", "-u")
	if err != nil {
		return err
	}

	mcs := findMappedCommands(all)
	for _, mc := range mcs {
		// gather the individual help text for each command
		gd, err := RunDoc(mc.FullName)
		if err != nil {
			panic(err)
		}
		lns := strings.SplitN(gd, "\n", 2)
		if len(lns) < 2 {
			continue
		}
		mc.Parameters = readLastParameters(lns[0])
		mc.HelpText = strings.TrimSpace(strings.Join(lns[1:], "\n"))
	}
	return json.NewEncoder(out).Encode(mcs)
}

func findMappedCommands(cmd string) MappedCommands {
	scn := bufio.NewScanner(bytes.NewBuffer([]byte(cmd)))
	cmds := []string{
		mainline.FuncNamePackage(mainline.MustAddCommand),
		mainline.FuncNamePackage(mainline.AddCommand),
	}

	mappings := make(MappedCommands)

	for scn.Scan() {
		ln := scn.Text()
		if ln == "" || strings.HasPrefix(ln, "//") {
			continue
		}
		for _, c := range cmds {
			if strings.Contains(ln, c) {
				pm := readLastParameters(ln)
				s := strings.Split(pm, ",")
				if len(s) != 2 {
					continue
				}
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
func readLastParameters(ln string) string {
	i := strings.LastIndex(ln, "(")
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
	args = append([]string{"doc"}, args...)
	c := exec.Command("go", args...)

	bufOut := bytes.NewBuffer(nil)
	bufErr := bytes.NewBuffer(nil)
	c.Stderr = bufErr
	c.Stdout = bufOut

	if err := c.Run(); err != nil {
		return "", fmt.Errorf("%s", bufErr.String())
	}
	return bufOut.String(), nil
}
