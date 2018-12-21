package query

import (
	"fmt"
	"net/http"
	"net/url"
)

type Query struct {
	scheme string
	host   string // includes port number
	method string

	header http.Header

	apiVersion   string
	namespace    string
	resourceType string
	resource     string

	query url.Values
}

func New(opts ...Opt) *Query {
	q := &Query{
		apiVersion: "/api/v1",
		namespace:  "default",

		scheme: "http",
		host:   "localhost",
		method: "GET",

		header: make(http.Header),
		query:  make(url.Values),
	}

	newQ := q.With(opts...)
	return newQ
}

func (q *Query) With(opts ...Opt) *Query {
	newQ := q

	for _, opt := range opts {
		newQ = opt(*newQ)
	}

	return newQ
}

func (q *Query) Request() *http.Request {
	req := &http.Request{
		Method: q.method,
		URL:    q.url(),
		Header: q.header,
	}

	return req
}

func (q *Query) url() *url.URL {
	return &url.URL{
		Scheme:   q.scheme,
		Host:     q.host,
		Path:     q.path(),
		RawQuery: q.query.Encode(),
	}
}

func (q *Query) path() string {
	return fmt.Sprintf(
		"%v/namespaces/%v/%v/%v",
		q.apiVersion,
		q.namespace,
		q.resourceType,
		q.resource,
	)
}
