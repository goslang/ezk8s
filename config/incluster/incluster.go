package incluster

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"github.com/goslang/ezk8s"
	"github.com/goslang/ezk8s/config"
	"github.com/goslang/ezk8s/query"
)

// Kubernetes publishes a service account and CA at a well known location in
// every Pod.
const (
	tokenFile  = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	rootCAFile = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
)

type clusterConfig struct{}

// New mimics Kubernetes' InCluster config behavior. Using this method will
// automatically pickup RBAC configured for the Pod the app is running in.
func New() config.Config {
	return &clusterConfig{}
}

// ClientOpts returns a list of ezk8s.Opts from the Kubernetes Pod's
// credentials.
func (cc *clusterConfig) ClientOpts() ([]ezk8s.Opt, error) {
	pool := x509.NewCertPool()
	err := addCerts(pool)
	if err != nil {
		return nil, err
	}

	host, err := getHost()
	if err != nil {
		return nil, err
	}

	token, err := getToken()
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
	}

	return []ezk8s.Opt{
		ezk8s.Transport(transport),
		ezk8s.QueryOpts(
			query.Host(host),
			query.AuthBearer(token),
		),
	}, nil
}

// Client creates a new ezk8s Client configured with the Kubernetes Pod's
// credentials. If the expected values cannot be found an error will be
// returned.
func (cc *clusterConfig) Client(opts ...ezk8s.Opt) (*ezk8s.Client, error) {
	defaults, err := cc.ClientOpts()
	if err != nil {
		return nil, err
	}

	cl := ezk8s.New(defaults...).With(opts...)
	return cl, nil
}

func getHost() (string, error) {
	host := os.Getenv("KUBERNETES_SERVICE_HOST")
	port := os.Getenv("KUBERNETES_SERVICE_PORT")

	if len(host) == 0 || len(port) == 0 {
		return "", fmt.Errorf(
			"Host and port not found, is this a Kubernetes cluster?",
		)
	}
	return "https://" + net.JoinHostPort(host, port), nil
}

func addCerts(pool *x509.CertPool) error {
	pem, err := ioutil.ReadFile(rootCAFile)
	if err != nil {
		return err
	}

	ok := pool.AppendCertsFromPEM(pem)
	if !ok {
		return fmt.Errorf("Invalid CA data for root CA file.")
	}
	return nil
}

func getToken() (string, error) {
	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return "", err
	}
	return string(token), nil
}
