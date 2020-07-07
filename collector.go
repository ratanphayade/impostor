package main

import "log"

type values map[string]string

func (v values) get(key string) string {
	if val, ok := v[key]; ok {
		return val
	}

	return ""
}

type collector struct {
	params    values
	resources values
	headers   values
	body      map[string]interface{}
}

func (c collector) getFromParam(key string) string {
	return c.params.get(key)
}

func (c collector) getFromResource(key string) string {
	return c.resources.get(key)
}

func (c collector) getFromHeader(key string) string {
	return c.headers.get(key)
}

func (c collector) getFromBody(key string) string {
	return get(c.body, key)
}

func (c collector) get(target string, modifier string) string {
	switch target {
	case TargetParams:
		return c.getFromParam(modifier)

	case TargetHeader:
		return c.getFromHeader(modifier)

	case TargetBody:
		return c.getFromBody(modifier)

	case TargetResource:
		return c.getFromResource(modifier)
	}

	log.Println("no matching target found")

	return ""
}
