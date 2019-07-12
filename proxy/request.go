/**
 * @Author: huangw1
 * @Date: 2019/7/11 11:00
 */

package proxy

import (
	"io"
	"net/url"
	"strings"
)

type Request struct {
	Method  string
	URL     *url.URL
	Query   url.Values
	Path    string
	Body    io.ReadCloser
	Params  map[string]string
	Headers map[string][]string
}

func (r *Request) GeneratePath(URLPattern string) {
	if len(r.Params) == 0 {
		r.Path = URLPattern
		return
	}
	for k, v := range r.Params {
		URLPattern = strings.Replace(URLPattern, "/{{."+k+"}}", "/"+v, -1)
	}
	r.Path = URLPattern
}

func (r *Request) Clone() Request {
	return Request{
		Method:  r.Method,
		URL:     r.URL,
		Query:   r.Query,
		Path:    r.Path,
		Body:    r.Body,
		Params:  r.Params,
		Headers: r.Headers,
	}
}
