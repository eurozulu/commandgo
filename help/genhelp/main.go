package main

import (
	"commandgo"
	"commandgo/help/genhelp/generate"
	"commandgo/help/genhelp/maps"
	"commandgo/help/genhelp/misc"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

const (
	defaultSrcPath     = "./..."
	defaultOutputName  = "cghelp.go"
	defaultPackageName = "main"
)

// ForceOverwrite when true will overwrite any existing maps file named in the Outname. (cghelp.go by default)
// When not present/false throws an os.Exists error if file already exists.
var ForceOverwrite bool

// Outname specifies the name of the generated go file containing help.  defaults to 'cghelp.go'.
// if the given name does not end with '.go', this will be added to the name.
var Outname = defaultOutputName

// PackageName specifies the package the generated file should be written into. Defaults to 'main'
// When specifying packages other than main, you must ensure that package is referenced somewhere in your application.
// Help relies on the 'init' function to register itself with the help system.  If unsure, leave it in main.``
var PackageName = defaultPackageName

func main() {
	cmd := commandgo.Commands{
		// one and only, default command
		"": MakeHelp,

		"-p": &PackageName,
		"-f": &ForceOverwrite,
		"-o": &Outname,
	}
	out, err := cmd.RunArgs()
	if err != nil {
		log.Fatalln(err)
	}
	for _, oi := range out {
		fmt.Println(oi)
	}
}

// MakeHelp generates a new help maps file containing the comments extracted from the files in the given maps directories.
// may specify one or more directory paths which contain the maps files to scan.  appending '...' to any path indicates to include all subdirectories of that path.
// maps locations should include all the code the command map maps into. i.e. the variables, structs, methods and func being mapped into.
func MakeHelp(src ...string) error {
	// check arguments, insert default src if empty and check Outname has .go
	if len(src) == 0 {
		src = []string{defaultSrcPath}
	}
	srcPaths, err := misc.collectSrcDirectories(src)
	if err != nil {
		return err
	}
	if len(srcPaths) == 0 {
		return fmt.Errorf("no maps files found in directorys: %s", strings.Join(srcPaths, " : "))
	}

	if !strings.HasSuffix(strings.ToLower(Outname), ".go") {
		Outname = strings.Join([]string{Outname, "go"}, ".")
	}

	pkgPath, err := misc.findPackage(PackageName, srcPaths)
	if err != nil {
		return err
	}
	outPath := path.Join(pkgPath, Outname)
	if err = checkOutPath(outPath); err != nil {
		return err
	}

	hgs, err := maps.makeHelp(srcPaths)
	if err != nil {
		return err
	}
	if len(hgs) == 0 {
		return fmt.Errorf("no maps files containing command maps were found")
	}

	by, err := generate.WriteTemplate(PackageName, hgs)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outPath, by, 0644)
}

// checkOutPath creates the path to the given path and checks if the filename exists
func checkOutPath(p string) error {
	if err := os.MkdirAll(path.Dir(p), 0755); err != nil {
		return err
	}
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		return nil
	}
	if err == nil {
		if !ForceOverwrite {
			return fmt.Errorf("%s already exists", p)
		}
	}
	return err
}
