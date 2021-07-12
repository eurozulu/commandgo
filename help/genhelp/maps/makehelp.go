package maps

import (
	"commandgo/help"
	"go/doc"
	"go/parser"
	"go/token"
)

const mapType = "commandgo.Commands"

type makeHelp struct {
	docs map[string]*doc.Package
}

// MakeHelp creates a group of Help subjects, scanning the given src paths for files using the commandgo.Commands map.
func (mh makeHelp) MakeHelp(srcPaths []string) ([]*help.HelpSubject, error) {
	var subjects []*help.HelpSubject

	// Find all instances of command maps in source
	var cmdMaps []PackageMaps
	for _, src := range srcPaths {
		pms, err := mh.findPackageMaps(src)
		if err != nil {
			return nil, err
		}
		if len(pms) == 0 {
			continue
		}
		cmdMaps = append(cmdMaps, pms...)
	}

	return subjects, nil
}

func (mh *makeHelp) findPackageMaps(srcPath string) ([]PackageMaps, error) {
	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, srcPath, nil, 0)
	if err != nil {
		return nil, err
	}
	var pms []PackageMaps
	for _, pkg := range pkgs {
		pm := NewPackageMaps(mapType, srcPath, pkg)
		if !pm.IsEmpty() {
			pms = append(pms, pm)
		}
	}
	return pms, nil
}

func (mh *makeHelp) addDocs(pm PackageMaps) []*help.HelpSubject {
	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, pm.SrcPath(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	for _, pd := range pkgs {
		d, ok := mh.docs[pkgName]
		if !ok {

		}
	}

}

func makeHelpItems(m map[string]string) []*help.HelpItem {
	// keep track of which values are mapped to and collect same mappings as aliases.
	vals := map[string]*help.HelpItem{}
	for k, v := range m {
		kk, ok := vals[v] // seen value before?
		if ok {
			// default "" takes precedence, otherwise, longest name wins
			if k == "" || len(k) > len(kk.Name) {
				s := k
				k = kk.Name
				kk.Name = s
			}
			kk.Aliases = append(kk.Aliases, k)
			continue
		}
		vals[k] = &help.HelpItem{
			Name:    k,
			Comment: v,
		}

	}
	return subjects
}
