package main

import (
	"bytes"
	"fmt"
	"github.com/eurozulu/mainline"
	"io"
	"net/http"
	"net/url"
	"os"
)

type SmallCurl struct {
	Header bool `flag:"header,i"`
}

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

	io.Copy(os.Stdout, r.Body)
	defer func() {
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
	}()
	return buf.String(), nil
}

func main() {
	mainline.AddCommand("get", SmallCurl.Get)
	mainline.RunCommandLine()
}
