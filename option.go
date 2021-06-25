package ezk8s

import (
	"net/http"

	"github.com/goslang/ezk8s/query"
)

// An Opt configures a single aspect of the ezk8s client.
type Opt func(c Client) *Client

// Transport configures the client to use the supplied http.RoundTripper when
// communicating with the Kubernetes API. This can be used to configure TLS.
func Transport(transport http.RoundTripper) Opt {
	return func(c Client) *Client {
		c.Client.Transport = transport
		return &c
	}
}

// QueryOpts sets default options to be used by all queries from this client.
// Typically this would be used to set the API Host, or similar options that
// should not change between requests.
func QueryOpts(opts ...query.Opt) Opt {
	return func(c Client) *Client {
		c.DefaultOpts = append(c.DefaultOpts, opts...)
		return &c
	}
}
