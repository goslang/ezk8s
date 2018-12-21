package query

// Opt returns a new Query with the provided configuration
type Opt func(Query) *Query

func Resource(resourceType, name string) Opt {
	return func(q Query) *Query {
		q.resourceType = resourceType
		q.resource = name
		return &q
	}
}

func ApiVersion(version string) Opt {
	return func(q Query) *Query {
		q.apiVersion = version
		return &q
	}
}

func Deployment(name string) Opt {
	resource := Resource("deployments", name)
	version := ApiVersion("/apis/apps/v1beta1")

	return func(q Query) *Query {
		return resource(*version(q))
	}
}

func Pod(name string) Opt {
	return Resource("pods", name)
}

func Label(name, value string) Opt {
	return func(q Query) *Query {
		q.query.Add("labelSelector", name+"="+value)
		return &q
	}
}

func Host(host string) Opt {
	return func(q Query) *Query {
		q.host = host
		return &q
	}
}

func Method(method string) Opt {
	return func(q Query) *Query {
		q.method = method
		return &q
	}
}

func AuthBearer(bearer string) Opt {
	return func(q Query) *Query {
		q.header.Add("Authorization", "Bearer "+bearer)
		return &q
	}
}
