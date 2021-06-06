package main

import (
	"encoding/json"
	"fmt"
	"io"
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

	streamPodEvents(cl)
}

func streamPodEvents(cl *ezk8s.Client) {
	reader, err := cl.Stream(
		query.Pod(""),
		query.Label("app.kubernetes.io/name", "hub-proxy"),
		query.Watch(""),
	)
	exitOnErr(err)
	defer reader.Close()

	decoder := json.NewDecoder(reader)
	for {
		var obj map[string]interface{}
		err := decoder.Decode(&obj)
		if err == io.EOF {
			break
		}
		exitOnErr(err)

		fmt.Println("type =", obj["type"])
	}
}

func exitOnErr(err error) {
	if err != nil {
		panic(err)
		fmt.Println(err)
		os.Exit(1)
	}
}
