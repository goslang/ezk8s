package main

import (
	"fmt"
	"os"

	"github.com/goslang/ezk8s"
	config "github.com/goslang/ezk8s/config/kube"
	"github.com/goslang/ezk8s/query"
)

func main() {
	conf, err := config.New("", "minikube")
	exitOnErr(err)

	cl, err := conf.Client()
	exitOnErr(err)

	getPersistentVolumes(cl)
}

func getPersistentVolumes(cl *ezk8s.Client) {
	res, err := cl.Query(
		query.PersistentVolume(""),
	)
	exitOnErr(err)

	var names []string
	err = res.Scan(
		query.Path{"$.items[:].metadata.name", &names},
	)
	exitOnErr(err)

	fmt.Println("persistent volume names:")
	for _, name := range names {
		fmt.Println(name)
	}
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
