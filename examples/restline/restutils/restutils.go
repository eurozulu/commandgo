// RestUtils is a simple example of an application having a command line interface applied to it.
// It performs simple http requests based on the given arguments.
package restutils

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	"io/ioutil"
	"net/http"
	"net/url"
)

// Verbose, when true, displays additional information about the operation.
var Verbose bool

// URLGet commands for requesting data from a remote location
type URLGet struct {
	// ShowHeaders, when true will display the response headers before the content
	// Also shows headers if Verbose is true.
	ShowHeaders bool

	// LocalFileRoot is optional and when set, defines the root location for local gets
	LocalFileRoot string
}

// URLPost commands for posting data to a remote site
type URLPost struct {
	// ContentType defines the format of the data being posted
	ContentType string

	// ShowHeaders displays the response headers before showing the response content
	ShowHeaders bool

	// LocalFilePermissions sets the file permissions of the local file when 'local' post is performed.
	LocalFilePermissions os.FileMode

	// LocalFileRoot is optional and when set, defines the root location for local posts
	LocalFileRoot string
}

// Get performs a HTTP GET operation on the given url, appending and given parameters to the url
func (g *URLGet) Get(u *url.URL, params ...string) (string, error) {
	r, err := http.Get(u.String())
	if err != nil {
		return "", err
	}

	out := bytes.NewBuffer(nil)
	if Verbose {
		out.WriteString(u.String())
		out.WriteString("\t")
		out.WriteString(fmt.Sprintf("%d\t%s", r.StatusCode, r.Status))
		out.WriteRune('\n')
	}
	if g.ShowHeaders || Verbose {
		out.WriteString("Headers:\n")
		for k, v := range r.Header {
			out.WriteString(k)
			out.WriteString(" = ")
			out.WriteString(strings.Join(v, ", "))
			out.WriteRune('\n')
		}
		out.WriteRune('\n')
	}

	by, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	out.Write(by)
	out.WriteRune('\n')

	return out.String(), nil
}

// GetLocal will retrieve a local file from the given path
func (g *URLGet) GetLocal(fn string) (string, error) {
	fn = path.Clean(fn)
	if g.LocalFileRoot != "" {
		if !strings.HasPrefix(fn, g.LocalFileRoot) {
			fn = path.Join(g.LocalFileRoot, fn)
		}
	}
	by, err := ioutil.ReadFile(fn)
	if err != nil {
		return "", err
	}
	return string(by), nil
}

// Post performs a http POST to the given URL, posting the given data
// u URL must be a valid http(s) URL
// data The data to post
// returns the status code and any response body or error
func (p *URLPost) Post(u *url.URL, data string) (int, string, error) {
	r, err := http.Post(u.String(), p.ContentType, strings.NewReader(data))
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	out := bytes.NewBuffer(nil)
	if Verbose {
		out.WriteString(u.String())
		out.WriteRune('\t')
		out.WriteString(fmt.Sprintf("%d\t%s", r.StatusCode, r.Status))
	}
	if p.ShowHeaders || Verbose {
		out.WriteString("Headers:\n")
		for k, v := range r.Header {
			out.WriteString(k)
			out.WriteString(" = ")
			out.WriteString(strings.Join(v, ", "))
			out.WriteRune('\n')
		}
		out.WriteRune('\n')
	}
	resp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return 0, "", err
	}
	return r.StatusCode, string(resp), nil
}

// PostLocal performs a local file save on the given data.
// fn File path where to write the given data
// data the data to write
// returns an error if failed to write
// use LocalFilePermissions to set a specific permission on the file, otherwise 0644 is used.
func (p *URLPost) PostLocal(fn string, data string) error {
	fn = path.Clean(fn)
	if p.LocalFileRoot != "" {
		if err := os.MkdirAll(fn, 0755); err != nil {
			return err
		}
		if !strings.HasPrefix(fn, p.LocalFileRoot) {
			fn = path.Join(p.LocalFileRoot, fn)
		}
	}
	if p.LocalFilePermissions == 0 {
		p.LocalFilePermissions = 0644
	}
	return ioutil.WriteFile(fn, []byte(data), p.LocalFilePermissions)
}
