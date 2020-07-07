package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/nqd/flat"
)

type (
	stringValues    map[string]string
	interfaceValues map[string]interface{}
)

func (v stringValues) get(key string) string {
	if val, ok := v[key]; ok {
		return val
	}

	return ""
}

func (v interfaceValues) get(key string) string {
	var result interface{}

	if val, ok := v[key]; ok {
		result = val
	} else {
		log.Println("key not found ", key)
		return ""
	}

	return fmt.Sprintf("%v", result)
}

type collector struct {
	params    stringValues
	resources stringValues
	headers   stringValues
	body      interfaceValues
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
	return c.body.get(key)
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

func collectRequestDetails(r *http.Request) collector {
	return collector{
		params:    collectParams(r),
		resources: collectResources(r),
		headers:   collectHeaders(r),
		body:      collectBody(r),
	}
}

func collectParams(r *http.Request) stringValues {
	data := stringValues{}

	for key, val := range r.URL.Query() {
		data[key] = val[0]
	}

	return data
}

func collectHeaders(r *http.Request) stringValues {
	data := stringValues{}

	for key, val := range r.Header {
		data[key] = val[0]
	}

	return data
}

func collectResources(r *http.Request) stringValues {
	data := stringValues{}

	for key, val := range mux.Vars(r) {
		data[key] = val
	}

	return data
}

func collectBody(r *http.Request) interfaceValues {
	var (
		request = interfaceValues{}
		data    interfaceValues
	)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println(err)
	}

	data, err := flat.Flatten(request, nil)
	if err != nil {
		log.Println(err)
	}

	return data
}
