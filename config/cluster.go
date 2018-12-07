package config

import (
	"crypto/x509"
	"errors"
	"io/ioutil"
)

var (
	ErrNoPEMData     = errors.New("No PEM data found in config for server CA.")
	ErrNoPEMFile     = errors.New("No PEM file for server CA.")
	ErrInvalidCAData = errors.New("Couldn't parse CA data for cluster.")
)

type cluster struct {
	Name        string
	ClusterData `yaml:"cluster"`
}

func (cl *cluster) loadServerCA() (*x509.CertPool, error) {
	pool := x509.NewCertPool()

	errs := [2]error{}
	errs[0] = cl.AddCertsFromData(pool)
	errs[1] = cl.AddCertsFromFile(pool)

	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}

	return pool, nil
}

func (cl *cluster) AddCertsFromData(pool *x509.CertPool) error {
	if cl.CertificateAuthorityData == "" {
		return ErrNoPEMData
	}

	ok := pool.AppendCertsFromPEM([]byte(cl.CertificateAuthorityData))
	if !ok {
		return ErrInvalidCAData
	}
	return nil
}

func (cl *cluster) AddCertsFromFile(pool *x509.CertPool) error {
	if cl.CertificateAuthority == "" {
		return ErrNoPEMFile
	}

	pem, err := ioutil.ReadFile(cl.CertificateAuthority)
	if err != nil {
		return err
	}

	ok := pool.AppendCertsFromPEM(pem)
	if !ok {
		return ErrInvalidCAData
	}
	return nil
}

type ClusterData struct {
	Server                   string
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
	CertificateAuthority     string `yaml:"certificate-authority"`
}

type clusters []cluster

func (cls clusters) Lookup(name string) (*cluster, bool) {
	for _, cluster := range cls {
		if cluster.Name == name {
			return &cluster, true
		}
	}

	return nil, false
}
