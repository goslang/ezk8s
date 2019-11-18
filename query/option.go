package query

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

// Opt returns a new Query with the provided configuration
type Opt func(Query) *Query

func Namespace(namespace string) Opt {
	return func(q Query) *Query {
		q.namespace = namespace
		return &q
	}
}

// Resource sets the resource type and name for the query.
func Resource(resourceType, name string) Opt {
	return func(q Query) *Query {
		q.resourceType = resourceType
		q.resource = name
		return &q
	}
}

// ApiVersion sets the api version that should be queried against.
func ApiVersion(version string) Opt {
	return func(q Query) *Query {
		q.apiVersion = version
		return &q
	}
}

// Deployment is a convenience method that sets the apiVersion, resourceType,
// and resource name for the Query. Passing an empty string will return all
// matching Deployments.
func Deployment(name string) Opt {
	resource := Resource("deployments", name)
	version := ApiVersion("/apis/apps/v1beta1")

	return func(q Query) *Query {
		return resource(*version(q))
	}
}

// Pod is a convenience method that sets resourceType and name for the Query.
// Passing an empty string will return all matching Pods.
func Pod(name string) Opt {
	return Resource("pods", name)
}

func Node(name string) Opt {
	resource := Resource("nodes", name)
	namespace := Namespace("")

	return func(q Query) *Query {
		return resource(*namespace(q))
	}
}

// Eviction is a convenience method for sending a pod Eviction to the
// Kubernetes API.
func Eviction(name string) Opt {
	resource := Resource("pods", name+"/eviction")
	method := Method("POST")
	reader := Json(map[string]interface{}{
		"apiVersion": "policy/v1beta1",
		"kind":       "Eviction",
		"metadata": map[string]interface{}{
			"name": name,
		},
	})

	return func(q Query) *Query {
		return reader(*resource(*method(q)))
	}
}

func Json(j interface{}) Opt {
	buf, _ := json.Marshal(j) // TODO: Ignoring error
	reader := bytes.NewReader(buf)

	return Body(ioutil.NopCloser(reader))
}

func Body(reader io.ReadCloser) Opt {
	return func(q Query) *Query {
		q.body = reader
		return &q
	}
}

// Label applies a labelSelector to the request.
func Label(name, value string) Opt {
	return func(q Query) *Query {
		q.query.Add("labelSelector", name+"="+value)
		return &q
	}
}

// Host sets the host name for the request. This will typically be set as a
// default by the ezk8s.Client.
func Host(host string) Opt {
	return func(q Query) *Query {
		q.host = host
		return &q
	}
}

// Sets the HTTP method for the requests. If not specified, the request will
// default to a GET.
func Method(method string) Opt {
	return func(q Query) *Query {
		q.method = method
		return &q
	}
}

// Sets the Bearer token for the request.
func AuthBearer(bearer string) Opt {
	return func(q Query) *Query {
		q.header.Add("Authorization", "Bearer "+bearer)
		return &q
	}
}
