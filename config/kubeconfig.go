package config

import (
	"errors"
)

var (
	ErrContextNotFound = errors.New("Context not found in kube config")
	ErrUserNotFound    = errors.New("User not found in kube config")
	ErrClusterNotFound = errors.New("Cluster not found in kube config")
)

type kubeConfig struct {
	CurrentContext string `yaml:"current-context"`

	Clusters clusters
	Users    users
	Contexts contexts
}

func (kc *kubeConfig) GetContext(name string) (*kubeContext, error) {
	ctx, ok := kc.Contexts.Lookup(name)
	if !ok {
		return nil, ErrContextNotFound
	}

	u, ok := kc.Users.Lookup(ctx.Context.User)
	if !ok {
		return nil, ErrUserNotFound
	}

	c, ok := kc.Clusters.Lookup(ctx.Context.Cluster)
	if !ok {
		return nil, ErrClusterNotFound
	}

	kubeCtx := &kubeContext{
		cluster: *c,
		user:    *u,
	}

	return kubeCtx, nil
}
