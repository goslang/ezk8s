package kube

import (
	"crypto/x509"
	"encoding/base64"
	"errors"
	"io/ioutil"
)

var (
	ErrNoPEMData     = errors.New("No PEM data found in config for server CA.")
	ErrNoPEMFile     = errors.New("No PEM file for server CA.")
	ErrInvalidCAData = errors.New("Couldn't parse CA data for cluster.")
)

type Cluster struct {
	Name        string
	ClusterData `yaml:"cluster"`
}

// loadServerCA returns the CA authorities for the server and an error if one was
// encountered. The final return will be true iff data was loaded, and false
// otherwise. This is necessary because it is possible to not load anything
// but still not fail, e.g. no CA was configured.
func (cl *Cluster) loadServerCA() (*x509.CertPool, error, bool) {
	pool := x509.NewCertPool()

	errs := [2]error{}
	errs[0] = cl.AddCertsFromData(pool)
	errs[1] = cl.AddCertsFromFile(pool)

	for idx, err := range errs {
		if err == nil {
			break
		} else if err != nil && idx == len(errs) {
			return nil, err, false
		}
	}

	return pool, nil, true
}

func (cl *Cluster) AddCertsFromData(pool *x509.CertPool) error {
	if cl.CertificateAuthorityData == "" {
		return ErrNoPEMData
	}

	pem, err := base64.StdEncoding.DecodeString(cl.CertificateAuthorityData)
	if err != nil {
		return err
	}

	ok := pool.AppendCertsFromPEM(pem)
	if !ok {
		return ErrInvalidCAData
	}
	return nil
}

func (cl *Cluster) AddCertsFromFile(pool *x509.CertPool) error {
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

type Clusters []Cluster

func (cls Clusters) Lookup(name string) (*Cluster, bool) {
	for _, cluster := range cls {
		if cluster.Name == name {
			return &cluster, true
		}
	}

	return nil, false
}
