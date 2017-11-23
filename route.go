package rip

import (
	"net/http"
	"path"
	"strings"
)

type Route struct {
	rip             *rip
	method          string
	doc             string
	target          interface{}
	path            []PathElement
	parameters      []Parameter
	queryParameters []*queryParameter
	handler         http.Handler
	errors          []error
}

func (r *Route) setMethod(method string) {
	if r.method != "" {
		r.errors = append(r.errors, &RouteMethodError{r, r.method, method})
	}
	r.method = method
}

func (r Route) SetMethod(m string) *Route {
	r.setMethod(m)
	return &r
}

func (r Route) Param(name string, doc string) *Route {
	qp := &queryParameter{
		name: name,
		doc:  doc,
	}
	r.parameters = append(r.parameters, qp)
	r.queryParameters = append(r.queryParameters, qp)
	return &r
}

func (r Route) Target(target interface{}) *Route {
	r.target = target
	return &r
}

func (r Route) Var(name string, doc string) *Route {
	varPath := &pathParameter{
		name: name,
		doc:  doc,
	}

	r.path = append(r.path, varPath)
	r.parameters = append(r.parameters, varPath)
	return &r
}

func (r Route) Path(pathElements ...string) *Route {
	for _, pe := range pathElements {
		pointer := fixedPath(pe)
		r.path = append(r.path, &pointer)
	}
	return &r
}

func (r Route) GET() *Route {
	r.setMethod("GET")
	return &r
}

func (r Route) POST() *Route {
	r.setMethod("POST")
	return &r
}

func (r Route) PUT() *Route {
	r.setMethod("PUT")
	return &r
}

func (r Route) DELETE() *Route {
	r.setMethod("DELETE")
	return &r
}

func (r *Route) Do(subcall func(*Route)) *Route {
	subcall(r)
	return r
}

func (r *Route) Handler(handler http.Handler, doc string) {
	if r.method == "" {
		r.errors = append(r.errors, &RouteMissingMethodError{r})
	}
	r.handler = handler
	r.doc = doc
	r.rip.routes = append(r.rip.routes, r)
}

func (r *Route) QueryTemplate() string {
	queryParts := []string{}
	for _, rp := range r.queryParameters {
		queryParts = append(queryParts, rp.Name())
	}

	if len(queryParts) > 0 {
		return "?" + strings.Join(queryParts, "&")
	}
	return ""
}

func (r *Route) Template() string {
	pathParts := []string{}
	for _, rp := range r.path {
		pathParts = append(pathParts, rp.Template())
	}
	template := "/" + path.Join(pathParts...)

	return template
}

func (r *Route) Method() string {
	return r.method
}
