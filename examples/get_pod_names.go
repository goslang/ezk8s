package main

import (
	"fmt"
	"os"

	"github.com/goslang/ezk8s"
	"github.com/goslang/ezk8s/config"
	"github.com/goslang/ezk8s/query"
)

func main() {
	conf, err := config.LoadFromKubeConfig("", "minikube")
	exitOnErr(err)

	cl := conf.Client()
	getPodNames(cl)
}

func getPodNames(cl *ezk8s.Client) {
	res, err := cl.Query(
		query.Pod(""),
		query.Label("app", "nginx"),
	)
	exitOnErr(err)

	var names []string
	err = res.Scan(
		query.Path{"$.items[:].metadata.name", &names},
	)
	exitOnErr(err)

	fmt.Println("pod names:")
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
