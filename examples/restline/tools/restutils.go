package tools

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Verbose, when true, displays additional information about the operation.
var Verbose bool

// URLGet is the command for requesting data from a remote location
type URLGet struct {
	// ShowHeaders, when true will display the response headers before the content
	// Also shows headers if Verbose is true.
	ShowHeaders bool
}

// URLPost manages posting data to a remote site
type URLPost struct {
	// ContentType defines the format of the data being posted
	ContentType string

	// Show the response headers
	ShowHeaders bool

	LocalFilePermissions os.FileMode
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

func (g *URLGet) GetLocal(p string) (string, error) {
	by, err := ioutil.ReadFile(p)
	if err != nil {
		return "", err
	}
	return string(by), nil
}

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

func (p *URLPost) PostLocal(fn string, data []byte) error {
	if p.LocalFilePermissions == 0 {
		p.LocalFilePermissions = 0644
	}
	return ioutil.WriteFile(fn, data, p.LocalFilePermissions)
}
