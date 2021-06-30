package kube

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os/exec"
	"time"
)

type Users []User

type User struct {
	Name     string
	UserData `yaml:"user"`
}

type UserData struct {
	ClientCertificate string `yaml:"client-certificate"`
	ClientKey         string `yaml:"client-key"`

	ClientCertificateData string `yaml:"client-certificate-data"`
	ClientKeyData         string `yaml:"client-key-data"`

	Exec *UserExec
}

type UserExec struct {
	Command            string
	Args               []string
	Env                map[string]string
	ProvideClusterInfo bool
}

// ExecCredential is the expected format returned by executing a "UserExec".
type ExecCredential struct {
	Kind       string
	ApiVersion string
	Status     struct {
		Token                 string
		ExpirationTimestamp   time.Time
		ClientCertificateData string
		ClientKeyData         string
	}
}

// ExecTripper is an http.RoundTripper that will inject the credentials
// returned by running "exec" into the request. The request will then be
// forwarded to "next".
type ExecTripper struct {
	exec UserExec
	next http.RoundTripper

	expiration time.Time
	creds      ExecCredential
}

func (et *ExecTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if et.expiration.IsZero() || time.Now().After(et.expiration.Add(time.Minute*-1)) {
		if err := et.load(); err != nil {
			return nil, err
		}
	}

	r.Header["Authorization"] = []string{"bearer " + et.creds.Status.Token}
	return et.next.RoundTrip(r)
}

func (et *ExecTripper) load() error {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := exec.Command(et.exec.Command, et.exec.Args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		// TODO: Return stderr output.
		return err
	}

	var creds ExecCredential
	if err := json.NewDecoder(stdout).Decode(&creds); err != nil {
		return err
	}

	et.creds = creds
	return nil
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
	cert, err := loadCertificateData(u.ClientCertificateData, u.ClientKeyData)
	return cert, err, err == nil
}

func (us Users) Lookup(name string) (*User, bool) {
	for _, u := range us {
		if u.Name == name {
			return &u, true
		}
	}
	return nil, false
}

func loadCertificateData(cert, key string) (tls.Certificate, error) {
	certData, err := base64.StdEncoding.DecodeString(cert)
	if err != nil {
		return tls.Certificate{}, err
	}

	keyData, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return tls.Certificate{}, err
	}

	return tls.X509KeyPair(certData, keyData)
}
