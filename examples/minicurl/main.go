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

// minicurl is a silly example of how small you can make a functional tool.
// Although crude, using a few lines of code, we can get a http response header with one command.
// get <url>
//
// This example shows how base functions can be used, on packages other than your own application.
package main

import (
	"github.com/eurozulu/mainline"
	"github.com/eurozulu/mainline/examples/minicurl/_help"
	"net/http"
)

func main() {
	mainline.MustAddCommand("get", http.Get)
	mainline.MustAddCommand("help", _help.HelpCommand)
	mainline.RunCommandLine()
}
