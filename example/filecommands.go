package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type FileCommands struct {
	Root      string
	Verbose   bool `flag:"v,verbose"`
	Recursive bool `flag:"recursive,r"`
}

func (cl FileCommands) ListFiles(p string) error {
	if cl.Root == "" {
		cl.Root, _ = os.Getwd()
	}
	fp := path.Join(cl.Root, p)
	fmt.Printf("listing files for %s ", fp)
	fmt.Println()
	names, err := cl.listPath(fp, true, false)
	if err != nil {
		return err
	}
	for _, n := range names {
		fmt.Println(n)
	}
	fmt.Println()
	return nil
}

func (cl FileCommands) ListDirectory(p string) error {
	if cl.Root == "" {
		cl.Root, _ = os.Getwd()
	}

	fmt.Printf("listing directories for %s ", p)
	if cl.Recursive {
		fmt.Print("recursively")
	}
	fmt.Println()

	names, err := cl.listPath(path.Join(cl.Root, p), false, true)
	if err != nil {
		return err
	}
	for _, n := range names {
		fmt.Println(n)
	}
	fmt.Println()
	return nil
}

func (cl FileCommands) listPath(p string, f, d bool) ([]string, error) {
	if !f && !d {
		return nil, nil
	}
	fis, err := ioutil.ReadDir(p)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, fi := range fis {
		n := []string{fi.Name()}

		if fi.IsDir() {
			if cl.Recursive {
				s, err := cl.listPath(path.Join(p, fi.Name()), f, d)
				if err != nil {
					return nil, err
				}
				n = append(n, s...)
			}
			if !d {
				continue
			}
			if cl.Verbose {
				n[0] = fmt.Sprintf("%s\t%v\t%v", n[0], fi.ModTime(), fi.Mode())
			}
		} else {
			if !f {
				continue
			}
			if cl.Verbose {
				n[0] = fmt.Sprintf("%s\t%v\t%v\t%d", n[0], fi.ModTime(), fi.Mode(), fi.Size())
			}
		}
		names = append(names, n...)
	}
	return names, nil
}
