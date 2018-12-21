package ezk8s

import (
	"encoding/json"
	"net/http"

	"github.com/goslang/ezk8s/query"
)

type Client struct {
	http.Client

	DefaultOpts []query.Opt
}

func New(opts ...Opt) *Client {
	cl := &Client{}

	return cl.With(opts...)
}

func (cl *Client) Query(opts ...query.Opt) (*query.Result, error) {
	q := cl.applyDefaults(
		query.New(opts...),
	)

	result := query.NewResult()
	req := q.Request()
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

func (cl *Client) applyDefaults(q *query.Query) *query.Query {
	return q.With(cl.DefaultOpts...)
}

func (cl *Client) With(opts ...Opt) (newCl *Client) {
	newCl = cl
	for _, opt := range opts {
		newCl = opt(*newCl)
	}
	return
}
