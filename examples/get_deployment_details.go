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

	getDeploymentDetails(cl)
}

func getDeploymentDetails(cl *ezk8s.Client) {
	res, err := cl.Query(
		query.Deployment("nginx-deployment"),
	)
	exitOnErr(err)

	var resourceVersion string
	var generation float64
	err = res.Scan(
		query.Path{"$.metadata.resourceVersion", &resourceVersion},
		query.Path{"$.metadata.generation", &generation},
	)
	exitOnErr(err)

	fmt.Println("Deployment details")
	fmt.Printf("generation = %v\n", generation)
	fmt.Printf("resourceVersion = %v\n", resourceVersion)
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
