package config

import (
	"crypto/tls"
)

// context represents the raw context data, before the relevant user/cluster
// information has been looked up.
type context struct {
	Name    string
	Context ContextData
}

type ContextData struct {
	Cluster string
	User    string `yaml:"user"`
}

type contexts []context

func (ctxs contexts) Lookup(name string) (*context, bool) {
	for _, ctx := range ctxs {
		if name == ctx.Name {
			return &ctx, true
		}
	}

	return nil, false
}

// kubeContext represents a fully parsed context object from a kube/config
// file.
type kubeContext struct {
	cluster cluster
	user    user
}

func (kc *kubeContext) loadTlsConfig() (*tls.Config, error) {
	tlsConf := tls.Config{}

	clientCert, err, didLoad := kc.user.loadClientTls()
	if err != nil {
		return nil, err
	} else if didLoad {
		tlsConf.Certificates = []tls.Certificate{clientCert}
	}

	cas, err, didLoad := kc.cluster.loadServerCA()
	if err != nil {
		return nil, err
	} else if didLoad {
		tlsConf.RootCAs = cas
	}

	tlsConf.BuildNameToCertificate()
	return &tlsConf, nil
}
