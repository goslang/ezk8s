package kube

import (
	"bytes"
	"encoding/json"
	//"fmt"
	"net/http"
	"os/exec"
	"time"
)

// ExecTripper is an http.RoundTripper that will inject the credentials
// returned by running "exec" into the request. The request will then be
// forwarded to "next".
type ExecTripper struct {
	exec UserExec
	next http.RoundTripper

	gate  chan bool
	creds ExecCredential
}

// ExecCredential is the expected format returned by executing a "UserExec".
type ExecCredential struct {
	Kind       string
	ApiVersion string
	Status     struct {
		Token               string
		ExpirationTimestamp time.Time

		// NOTE: These would be needed to support TLS certs from a command.
		//ClientCertificateData string
		//ClientKeyData         string
	}
}

func NewExecTripper(exec UserExec, next http.RoundTripper) *ExecTripper {
	tripper := &ExecTripper{
		exec: exec,
		next: next,
		gate: make(chan bool, 1),
	}

	tripper.gate <- true
	return tripper
}

func (et *ExecTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if _, ok := <-et.gate; ok {
		// A successful read means that the credentials need to be loaded.
		// All other goroutines will remain blocked on the read until the
		// channel is closed, after the credentials have been updated.
		// A timer will be kicked off to reinitialize the buffered channel,
		// triggering the whole process again.

		func() {
			defer func() { close(et.gate) }()
			et.load()
		}()

		go func() {
			exp := et.creds.Status.ExpirationTimestamp
			timeout := exp.Sub(time.Now().Add(time.Minute * -1))
			<-time.After(timeout)
			et.gate = make(chan bool, 1)
			et.gate <- true
		}()
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
