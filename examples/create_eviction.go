package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/goslang/ezk8s"
	config "github.com/goslang/ezk8s/config/kube"
	"github.com/goslang/ezk8s/query"
)

func main() {
	conf, err := config.New("", "proxy")
	exitOnErr(err)

	podName := flag.String("pod", "", "The pod to evict")
	flag.Parse()

	cl, err := conf.Client()
	exitOnErr(err)

	evictPod(*podName, cl)
}

func evictPod(podName string, cl *ezk8s.Client) {
	fmt.Println("Evicting", podName)

	_, err := cl.Query(
		query.Eviction(podName),
	)
	exitOnErr(err)

	fmt.Printf("Evicted %v\n", podName)
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
