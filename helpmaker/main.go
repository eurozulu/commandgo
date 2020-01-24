package main

import (
	"github.com/go-bindata/go-bindata"
	"os"
	"path"
	"strings"
)

const helpdir = "_help"
const helpJSON = "help.json"
const helpGo = "help.go"

func main() {
	//by, err := ioutil.ReadFile("help.txt")
	var p string
	if len(os.Args) < 2 {
		p, _ = os.Getwd()
	} else {
		p = os.Args[1]
	}

	hd := path.Join(p, helpdir)
	if err := os.MkdirAll(hd, 0644); err != nil {
		panic(err)
	}

	jfp := path.Join(hd, helpJSON)

	fo, err := os.OpenFile(jfp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err := os.MkdirAll(hd, 0644); err != nil {
		panic(err)
	}

	err = ParseCommands(p, fo)
	errc := fo.Close()

	if err != nil {
		panic(err)
	}
	if errc != nil {
		panic(err)
	}

	if err := buildResource(hd, helpGo); err != nil {
		panic(err)
	}
}

func buildResource(pkg string, outName string) error {
	in := bindata.InputConfig{
		Path:      strings.SplitN(pkg, "/...", 2)[0], // strip off trailing ... if present
		Recursive: true,
	}

	cfg := bindata.Config{
		Package:    pkg,
		Input:      []bindata.InputConfig{in},
		Output:     outName,
		NoMetadata: true,
	}

	return bindata.Translate(&cfg)
}
