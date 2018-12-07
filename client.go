package ezk8s

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tma1/ezk8s/query"
)

type Client struct {
	http.Client

	DefaultOpts []query.Opt
}

func New(defaults ...query.Opt) *Client {
	return &Client{
		DefaultOpts: defaults,
	}
}

func (cl *Client) Query(opts ...query.Opt) (*query.Result, error) {
	q := cl.applyDefaults(
		query.New(opts...),
	)

	result := query.NewResult()
	req := q.Request()
	fmt.Println(req.Header)
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
