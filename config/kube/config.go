package kube

import (
	"os"
	osUser "os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/goslang/ezk8s/config"
)

func New(path, contextName string) (config.Config, error) {
	kubePath := getKubeConfigPath(path)
	file, err := os.Open(kubePath)
	if err != nil {
		return nil, err
	}

	k8Conf := kubeConfig{}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&k8Conf); err != nil {
		return nil, err
	}

	k8Ctx, err := k8Conf.GetContext(contextName)
	if err != nil {
		return nil, err
	}

	return k8Ctx, nil
}

func getKubeConfigPath(path string) string {
	if path != "" {
		return path
	}

	if usr, err := osUser.Current(); err == nil {
		return filepath.Join(usr.HomeDir, ".kube/config")
	}

	return ".kube/config"
}
