package config

import (
	"crypto/tls"
)

type user struct {
	Name     string
	UserData `yaml:"user"`
}

// loadClientTls returns the Certificate and an error if encountered
// In the case that no data was configured, error will be nil, but the
// certificate will still be a zero value. If this occurs, the final bool will
// return false, indicating that nothing was loaded, even though there was no
// error.
func (u *user) loadClientTls() (tls.Certificate, error, bool) {
	if u.hasCertData() {
		return u.loadCertificateData()
	}

	if u.hasCertFiles() {
		return u.loadCertificateFiles()
	}

	return tls.Certificate{}, nil, false
}

func (u *user) hasCertData() bool {
	return u.ClientCertificateData != "" && u.ClientKeyData != ""
}

func (u *user) hasCertFiles() bool {
	return u.ClientCertificate != "" && u.ClientKey != ""
}

func (u *user) loadCertificateFiles() (tls.Certificate, error, bool) {
	cert, err := tls.LoadX509KeyPair(
		u.ClientCertificate,
		u.ClientKey,
	)

	if err != nil {
		return cert, err, false
	}
	return cert, nil, true
}

func (u *user) loadCertificateData() (tls.Certificate, error, bool) {
	cert, err := tls.LoadX509KeyPair(
		u.ClientCertificate,
		u.ClientKey,
	)

	if err != nil {
		return cert, err, false
	}
	return cert, nil, true
}

type UserData struct {
	ClientCertificate string `yaml:"client-certificate"`
	ClientKey         string `yaml:"client-key"`

	ClientCertificateData string `yaml:"client-certificate-data"`
	ClientKeyData         string `yaml:"client-key-data"`
}

type users []user

func (us users) Lookup(name string) (*user, bool) {
	for _, u := range us {
		if u.Name == name {
			return &u, true
		}
	}
	return nil, false
}
