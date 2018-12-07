package config

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

const kubePath = "/Users/andrewknapp/.kube/config"

type Config struct {
	transport *http.Transport
}

func LoadFromKubeConfig() (*Config, error) {
	file, err := os.Open(kubePath)
	if err != nil {
		return nil, err
	}

	k8Conf := kubeConfig{}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&k8Conf); err != nil {
		return nil, err
	}

	fmt.Println(k8Conf.GetContext("minikube"))
	return nil, nil
}

func buildTlsTransport(pemKey, pemCert []byte) (*http.Transport, error) {
	cert, err := tls.X509KeyPair(pemKey, pemCert)
	if err != nil {
		return nil, err
	}

	tlsConfig := tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	tlsConfig.BuildNameToCertificate()

	transport := http.Transport{
		TLSClientConfig: &tlsConfig,
	}

	return &transport, nil
}
