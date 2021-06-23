package main

import (
	"flag"
	"fmt"
	"os"

	config "github.com/goslang/ezk8s/config/kube"
	"github.com/goslang/ezk8s/query"
)

func main() {
	context := flag.String("context", "minikube", "The config context to use")
	name := flag.String("node", "", "The node to cordon")
	enabled := flag.Bool("enable", false, "Disable the node")
	flag.Parse()

	conf, err := config.New("", *context)

	cl, err := conf.Client()
	exitOnErr(err)

	node := make(map[string]interface{})
	err = cl.Query(query.Node(*name)).Decode(&node)
	exitOnErr(err)

	setNodeState(*enabled, node)
	err = cl.Query(
		query.Node(*name),
		query.Method("PUT"),
		query.Json(node),
	).Error()
	exitOnErr(err)
}

func setNodeState(enabled bool, node map[string]interface{}) {
	spec, ok := node["spec"].(map[string]interface{})
	if !ok {
		spec = make(map[string]interface{})
		node["spec"] = spec
	}
	spec["unschedulable"] = !enabled
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
