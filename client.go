package ezk8s

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/goslang/ezk8s/query"
)

// Client is responsible for handling communication details with the
// Kubernetes API.
type Client struct {
	http.Client

	DefaultOpts []query.Opt
}

// New creates a new ezk8s.Client and applies the supplied options.
func New(opts ...Opt) *Client {
	cl := &Client{}

	return cl.With(opts...)
}

// Query sends a request to the Kubernetes API and returns the result. If an
// error occurred during the request, calling any method on the Result will
// return that error.
func (cl *Client) Query(opts ...query.Opt) query.Result {
	q := cl.applyDefaults(
		query.New(opts...),
	)

	req, err := q.Request()
	if err != nil {
		return query.NewErrorResult(err)
	}

	response, err := cl.Client.Do(req)
	if err != nil {
		return query.NewErrorResult(err)
	}

	if response.StatusCode >= 300 || response.StatusCode < 200 {
		defer response.Body.Close()

		buf, _ := ioutil.ReadAll(response.Body)
		return query.NewErrorResult(fmt.Errorf(
			"Error Response code %v\nresponse body: %s",
			response.StatusCode,
			buf,
		))
	}

	return query.NewDecodeResult(response.Body)
}

// With creates a new client after applying the supplied options.
func (cl *Client) With(opts ...Opt) (newCl *Client) {
	newCl = cl
	for _, opt := range opts {
		newCl = opt(*newCl)
	}
	return
}

func (cl *Client) applyDefaults(q *query.Query) *query.Query {
	return q.With(cl.DefaultOpts...)
}
