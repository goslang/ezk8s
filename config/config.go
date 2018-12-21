package config

import (
	"crypto/tls"
	"net/http"
	"os"
	osUser "os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"

	client "github.com/goslang/ezk8s"
	"github.com/goslang/ezk8s/query"
)

type Config struct {
	TlsConfig *tls.Config

	Host string
}

func LoadFromKubeConfig(path, contextName string) (*Config, error) {
	conf := &Config{}

	kubePath := getKubeConfigPath(path)
	file, err := os.Open(kubePath)
	if err != nil {
		return nil, err
	}

	k8Conf := kubeConfig{}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&k8Conf); err != nil {
		return nil, err
	}

	k8Ctx, err := k8Conf.GetContext(contextName)
	if err != nil {
		return nil, err
	}

	conf.TlsConfig, err = k8Ctx.loadTlsConfig()
	if err != nil {
		return nil, err
	}

	conf.Host = k8Ctx.cluster.ClusterData.Server

	return conf, nil
}

func getKubeConfigPath(path string) string {
	if path != "" {
		return path
	}

	if usr, err := osUser.Current(); err == nil {
		return filepath.Join(usr.HomeDir, ".kube/config")
	}

	return ".kube/config"
}

func (c *Config) ClientOpts() []client.Opt {
	queryOpts := []query.Opt{
		query.Host(c.Host),
	}

	return []client.Opt{
		client.Transport(buildTlsTransport(c.TlsConfig)),
		client.QueryOpts(queryOpts...),
	}
}

func (c *Config) Client(opts ...client.Opt) *client.Client {
	return client.New(opts...).With(c.ClientOpts()...)
}

func buildTlsTransport(tlsConfig *tls.Config) *http.Transport {
	return &http.Transport{
		TLSClientConfig: tlsConfig,
	}
}
