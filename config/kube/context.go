package kube

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/goslang/ezk8s"
	"github.com/goslang/ezk8s/query"
)

// context represents the raw context data, before the relevant user/cluster
// information has been looked up.
type Context struct {
	Name    string
	Context ContextData
}

type ContextData struct {
	Cluster string
	User    string `yaml:"user"`
}

type Contexts []Context

func (ctxs Contexts) Lookup(name string) (*Context, bool) {
	for _, ctx := range ctxs {
		if name == ctx.Name {
			return &ctx, true
		}
	}

	return nil, false
}

// KubeContext manages configuration from a .kube/config context and
// implements config.Config.
type KubeContext struct {
	Cluster Cluster
	User    User
}

// ClientOpts returns the list of options that should be past to ezk8s.New to
// correctly configure the client.
func (kc *KubeContext) ClientOpts() (opts []ezk8s.Opt, err error) {
	defer recoverPanic(&err)

	// Build the default query.Opts
	queryOpts := []query.Opt{
		query.Host(kc.Cluster.ClusterData.Server),
	}

	transport := kc.buildTlsTransport()
	if kc.User.UserData.Exec != nil {
		transport = NewExecTripper(*kc.User.UserData.Exec, transport)
	}

	// Build the client.Opts
	opts = []ezk8s.Opt{
		ezk8s.Transport(transport),
		ezk8s.QueryOpts(queryOpts...),
	}

	return
}

// Client builds a new ezk8s client with the options acquired from ClientOpts.
func (kc *KubeContext) Client(opts ...ezk8s.Opt) (*ezk8s.Client, error) {
	defaults, err := kc.ClientOpts()
	if err != nil {
		return nil, err
	}

	cl := ezk8s.New(defaults...).With(opts...)
	return cl, nil
}

// buildTlsTransport builds an http.Transport that includes TLS details.
// Panics on error.
func (kc *KubeContext) buildTlsTransport() http.RoundTripper {
	return &http.Transport{
		TLSClientConfig: kc.loadTlsConfig(),
	}
}

// Builds a tls.Config that includes both server and client TLS details.
// Panics on error.
func (kc *KubeContext) loadTlsConfig() *tls.Config {
	tlsConf := tls.Config{}

	clientCert, err, didLoad := kc.User.loadClientTls()
	if err != nil {
		panic(err)
	} else if didLoad {
		tlsConf.Certificates = []tls.Certificate{clientCert}
	}

	cas, err, didLoad := kc.Cluster.loadServerCA()
	if err != nil {
		panic(err)
	} else if didLoad {
		tlsConf.RootCAs = cas
	}

	tlsConf.BuildNameToCertificate()
	return &tlsConf
}

func recoverPanic(target *error) {
	if err := recover(); err != nil {
		switch e := err.(type) {
		case error:
			*target = e
		default:
			*target = fmt.Errorf("%v", e)
		}
	}
}
