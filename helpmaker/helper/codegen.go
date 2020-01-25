// Copyright 2020 Rob Gilham
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

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
