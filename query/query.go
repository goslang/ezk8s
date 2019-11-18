package query

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Query represents a single request to the Kubernetes API that.
type Query struct {
	method string
	host   string

	header http.Header

	apiVersion   string
	namespace    string
	resourceType string
	resource     string

	body io.ReadCloser

	query url.Values
}

// New returns a new query configured with the supplied options. It also
// attempts to use sane defaults.
func New(opts ...Opt) *Query {
	q := &Query{
		apiVersion: "/api/v1",
		namespace:  "default",

		method: "GET",
		host:   "http://localhost",

		header: make(http.Header),
		query:  make(url.Values),
	}

	newQ := q.With(opts...)
	return newQ
}

// With applies the options to a new Query based off of the old one.
func (q *Query) With(opts ...Opt) *Query {
	newQ := q

	for _, opt := range opts {
		newQ = opt(*newQ)
	}

	return newQ
}

// Request returns the HTTP representation of the Query, suitable for use by
// an http.Client.
func (q *Query) Request() (*http.Request, error) {
	reqUrl, err := q.url()
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: q.method,
		URL:    reqUrl,
		Header: q.header,
		Body:   q.body,
	}

	return req, nil
}

func (q *Query) url() (*url.URL, error) {
	fullUrl, err := url.Parse(q.host)
	if err != nil {
		return nil, err
	}

	fullUrl.Path = q.path()
	fullUrl.RawQuery = q.query.Encode()

	return fullUrl, nil
}

func (q *Query) path() string {
	parts := []string{q.apiVersion}

	push := func(strs ...string) {
		for _, s := range strs {
			if s == "" {
				return
			}
		}

		parts = append(parts, strs...)
	}

	push("namespaces", q.namespace)
	push(q.resourceType)
	push(q.resource)

	return strings.Join(parts, "/")
}
