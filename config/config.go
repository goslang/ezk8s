package config

import (
	"github.com/goslang/ezk8s"
)

type Config interface {
	Client(...ezk8s.Opt) (*ezk8s.Client, error)
	ClientOpts() ([]ezk8s.Opt, error)
}
