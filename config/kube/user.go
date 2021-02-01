package kube

import (
	"crypto/tls"
)

type User struct {
	Name     string
	UserData `yaml:"user"`
}

// loadClientTls returns the Certificate and an error if encountered
// In the case that no data was configured, error will be nil, but the
// certificate will still be a zero value. If this occurs, the final bool will
// return false, indicating that nothing was loaded, even though there was no
// error.
func (u *User) loadClientTls() (tls.Certificate, error, bool) {
	if u.hasCertData() {
		return u.loadCertificateData()
	}

	if u.hasCertFiles() {
		return u.loadCertificateFiles()
	}

	return tls.Certificate{}, nil, false
}

func (u *User) hasCertData() bool {
	return u.ClientCertificateData != "" && u.ClientKeyData != ""
}

func (u *User) hasCertFiles() bool {
	return u.ClientCertificate != "" && u.ClientKey != ""
}

func (u *User) loadCertificateFiles() (tls.Certificate, error, bool) {
	cert, err := tls.LoadX509KeyPair(
		u.ClientCertificate,
		u.ClientKey,
	)

	if err != nil {
		return cert, err, false
	}
	return cert, nil, true
}

func (u *User) loadCertificateData() (tls.Certificate, error, bool) {
	cert, err := tls.LoadX509KeyPair(
		u.ClientCertificateData,
		u.ClientKeyData,
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

type Users []User

func (us Users) Lookup(name string) (*User, bool) {
	for _, u := range us {
		if u.Name == name {
			return &u, true
		}
	}
	return nil, false
}
