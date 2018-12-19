package ezk8s

import (
	"net/http"

	"github.com/tma1/ezk8s/query"
)

type Opt func(c Client) *Client

func Transport(transport http.RoundTripper) Opt {
	return func(c Client) *Client {
		c.Client.Transport = transport
		return &c
	}
}

func QueryOpts(opts ...query.Opt) Opt {
	return func(c Client) *Client {
		c.DefaultOpts = opts
		return &c
	}
}
