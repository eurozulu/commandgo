// smallcurl is another triviol example which builds on the mini curl.
// Still keeping it simple, it adds a struct to handle the response.
// With the struct it can use flags, so adds a flag to get the headers as well as the body
package main

import (
	"bytes"
	"fmt"
	"github.com/eurozulu/mainline"
	"github.com/eurozulu/mainline/examples/smallcurl/_help"
	"io"
	"net/http"
	"net/url"
	"os"
)

type SmallCurl struct {
	// Header, when true, displays the response headers at the end of the body
	Header bool `flag:"header,i"`
	// Nobody, when true, will NOT return the body stream.  Used with Header to get just headers, or on its own gets response code.
	Nobody bool
}

// Get performs a http GET on the given url, returning the BODY
func (sc SmallCurl) Get(u *url.URL) (string, error) {
	r, err := http.Get(u.String())
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(nil)
	if sc.Header {
		for k, v := range r.Header {
			buf.WriteString(fmt.Sprintf("%s = %v\n", k, v))
		}
		fmt.Println()
	}

	if !sc.Nobody {
		io.Copy(os.Stdout, r.Body)
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
	}()

	if buf.Len() == 0 {
		buf.WriteString(r.Status)
	}
	return buf.String(), nil
}

func main() {
	mainline.AddCommand("get", SmallCurl.Get)
	mainline.MustAddCommand("help", _help.HelpCommand)
	mainline.RunCommandLine()
}
