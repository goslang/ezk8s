package config

import (
	"crypto/tls"
)

type user struct {
	Name     string
	UserData `yaml:"user"`
}

func (u *user) loadClientTls() (tls.Certificate, error) {
	if u.hasCertData() {
		return tls.X509KeyPair(
			[]byte(u.ClientCertificateData),
			[]byte(u.ClientKeyData),
		)
	}

	if u.hasCertFiles() {
		return tls.LoadX509KeyPair(
			u.ClientCertificate,
			u.ClientKey,
		)
	}

	return tls.LoadX509KeyPair(
		u.ClientCertificate,
		u.ClientKey,
	)
}

func (u *user) hasCertData() bool {
	return u.ClientCertificateData != "" && u.ClientKeyData != ""
}

func (u *user) hasCertFiles() bool {
	return u.ClientCertificate != "" && u.ClientKey != ""
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
