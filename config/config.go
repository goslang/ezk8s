package config

import (
	"github.com/goslang/ezk8s"
)

// Config objects are used to load default options from a variety of sources,
// such as a kube/config.
type Config interface {
	Client(...ezk8s.Opt) (*ezk8s.Client, error)
	ClientOpts() ([]ezk8s.Opt, error)
}
