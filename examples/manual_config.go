package main

import (
	"fmt"
	"os"

	"github.com/goslang/ezk8s"
	"github.com/goslang/ezk8s/query"
)

func main() {
	cl := &ezk8s.Client{}

	res := cl.Query(
		query.Host("http://127.0.0.1:8001"),
		query.Pod(""),
	)

	var names []string
	err := res.Scan(query.Path{"$.items[0:-1].metadata.name", &names})
	exitOnErr(err)

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
