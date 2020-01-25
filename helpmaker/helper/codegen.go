package helper

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func GenCode(by []byte) error {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("package %s\n", helpdir))
	buf.WriteRune('\n')
	buf.WriteString(HelpImport)
	buf.WriteRune('\n')
	buf.WriteString(fmt.Sprintf("const helpData = `%s`", string(by)))
	buf.WriteString(HelpFunction)

	if err := os.MkdirAll(helpdir, 0755); err != nil {
		return err
	}
	outName := path.Join(helpdir, helpGo)
	return ioutil.WriteFile(outName, buf.Bytes(), 0644)
}
