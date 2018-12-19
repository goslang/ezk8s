package config

import (
	"crypto/tls"
	"net/http"
	"os"
	osUser "os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"

	client "github.com/tma1/ezk8s"
)

type Config struct {
	TlsConfig *tls.Config
}

func LoadFromKubeConfig(contextName string) (*Config, error) {
	kubePath := getKubeConfigPath()
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

	tlsConfig, err := k8Ctx.loadTlsConfig()
	if err != nil {
		return nil, err
	}

	conf := &Config{
		TlsConfig: tlsConfig,
	}

	return conf, nil
}

func getKubeConfigPath() string {
	if usr, err := osUser.Current(); err == nil {
		return filepath.Join(usr.HomeDir, ".kube/config")
	}
	return ""
}

func (c *Config) ClientOpts() []client.Opt {
	return []client.Opt{
		client.Transport(buildTlsTransport(c.TlsConfig)),
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
