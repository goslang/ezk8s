package query

import (
	"fmt"
	"net/http"
	"net/url"
)

type Query struct {
	method string
	host   string

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

		method: "GET",
		host:   "http://localhost",

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

func (q *Query) Request() (*http.Request, error) {
	reqUrl, err := q.url()
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: q.method,
		URL:    reqUrl,
		Header: q.header,
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
	return fmt.Sprintf(
		"%v/namespaces/%v/%v/%v",
		q.apiVersion,
		q.namespace,
		q.resourceType,
		q.resource,
	)
}
