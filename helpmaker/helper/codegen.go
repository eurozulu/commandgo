package helper

import (
	"bufio"
	"bytes"
	"github.com/go-bindata/go-bindata"
	"io/ioutil"
	"path"
	"strings"
)

func GenCode() error {
	outName := path.Join(helpdir, helpGo)
	in := bindata.InputConfig{
		Path:      strings.SplitN(helpdir, "/...", 2)[0], // strip off trailing ... if present
		Recursive: true,
	}

	cfg := bindata.Config{
		Package:    helpdir,
		Input:      []bindata.InputConfig{in},
		Output:     outName,
		NoMetadata: true,
	}

	err := bindata.Translate(&cfg)
	if err != nil {
		return err
	}

	gf, err := ioutil.ReadFile(outName)
	if err != nil {
		return err
	}
	gf = injectImport(HelpImport, gf)
	gf = injectFunction(HelpFunction, gf)
	return ioutil.WriteFile(outName, gf, 0644)
}

func injectFunction(fun string, by []byte) []byte {
	l := len(by)
	b := make([]byte, len(fun)+l)
	copy(b, by)
	copy(b[l:], []byte(fun))
	return b
}

func injectImport(pkg string, by []byte) []byte {
	ip := findPackage(by)
	if ip < 0 {
		panic("No package declaration found")
	}
	b := make([]byte, len(by)+len(HelpImport))
	copy(b, by[0:ip])
	copy(b[ip:], HelpImport)
	copy(b[ip+len(HelpImport):], by[ip:])
	return b
}

func findPackage(by []byte) int {
	scn := bufio.NewScanner(bytes.NewBuffer(by))
	i := 0
	for scn.Scan() {
		s := scn.Text()
		i += len(s) + 1
		if strings.HasPrefix(s, "package") {
			return i
		}
	}
	return -1
}
