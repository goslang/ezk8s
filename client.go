package ezk8s

import (
	"encoding/json"
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

// Query sends a request to the Kubernetes API and returns the result, or an
// error.
func (cl *Client) Query(opts ...query.Opt) (*query.Result, error) {
	q := cl.applyDefaults(
		query.New(opts...),
	)

	result := query.NewResult()
	req, err := q.Request()
	if err != nil {
		return nil, err
	}

	response, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&result.Data); err != nil {
		return nil, err
	}

	return result, nil
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
