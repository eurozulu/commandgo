package misc

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func inSlice(s string, ss []string) bool {
	for _, sz := range ss {
		if sz == s {
			return true
		}
	}
	return false
}

// findPackage locates the directory of containing the named package.
func findPackage(name string, srcPaths []string) (string, error) {
	names, err := collectSrcDirectories(srcPaths)
	if err != nil {
		return "", err
	}

	for _, dir := range names {
		fs := token.NewFileSet()
		pkgs, err := parser.ParseDir(fs, dir, nil, 0)
		if err != nil {
			return "", err
		}
		for _, pkg := range pkgs {
			if pkg.Name != name {
				continue
			}
			return dir, nil
		}
	}
	return "", fmt.Errorf("%s package could not be found", name)
}

// collectSrcDirectories colelcts the full path names of the given directories which contain .go files.
// if any of the given paths end with "/..." the sub directories of that path are also searched.
func collectSrcDirectories(paths []string) ([]string, error) {
	var found []string
	for _, p := range paths {
		r := strings.HasSuffix(p, "/...")
		if r {
			p = strings.TrimRight(p, ".")
		}
		fis, err := ioutil.ReadDir(p)
		if err != nil {
			return nil, err
		}

		if containsGoFile(fis) {
			found = append(found, p)
		}

		if !r {
			continue
		}
		sds, err := collectSrcDirectories(subDirectories(p, fis))
		found = append(found, sds...)
	}
	return found, nil
}

func containsGoFile(fis []os.FileInfo) bool {
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(fi.Name(), ".go") {
			return true
		}
	}
	return false
}

func subDirectories(root string, fis []os.FileInfo) []string {
	var found []string
	for _, fi := range fis {
		if !fi.IsDir() {
			continue
		}
		found = append(found, path.Join(root, fi.Name(), "..."))
	}
	return found
}
