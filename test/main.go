package main

import (
	"fmt"
	"os"

	k8s "github.com/tma1/ezk8s"
	"github.com/tma1/ezk8s/query"
)

func main() {
	cl := k8s.New(
		query.Host("127.0.0.1:8001"),
		query.Scheme("http"),
	)

	res, err := cl.Query(
		query.Deployment("nginx-deployment"),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var resourceVersion string
	var generation float64
	err = res.Scan(
		query.Path{"$.metadata.resourceVersion", &resourceVersion},
		query.Path{"$.metadata.generation", &generation},
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf(
		"generation = %v\nresourceVersion = %v\n",
		generation, resourceVersion,
	)
}
