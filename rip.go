package rip

import (
	"github.com/gorilla/mux"
	"net/http"
)

type rip struct {
	routes []*Route
}

func NewRIP() *rip {
	return &rip{}
}

func (rip *rip) Routes() []*Route {
	return rip.routes
}

func (rip *rip) routeErrors() []error {
	errors := []error{}
	for _, r := range rip.routes {
		errors = append(errors, r.errors...)
	}
	return errors
}

func (rip *rip) AsRoute() *Route {
	return &Route{rip: rip}
}

func (rip *rip) Param(name string, doc string) *Route {
	return Route{rip: rip}.Param(name, doc)
}

func (rip *rip) Var(name string, doc string) *Route {
	return Route{rip: rip}.Var(name, doc)
}

func (rip *rip) Path(pathElements ...string) *Route {
	return Route{rip: rip}.Path(pathElements...)
}

func (rip *rip) GET() *Route {
	return Route{rip: rip}.GET()
}

func (rip *rip) POST() *Route {
	return Route{rip: rip}.POST()
}

func (rip *rip) PUT() *Route {
	return Route{rip: rip}.PUT()
}

func (rip *rip) DELETE() *Route {
	return Route{rip: rip}.DELETE()
}

func (rip *rip) Do(subcall func(*Route)) *Route {
	r := &Route{rip: rip}
	return r.Do(subcall)
}

func (rip *rip) RootHandler() (http.Handler, error) {
	errors := rip.routeErrors()
	if len(errors) > 0 {
		return nil, errors[0]
	}

	return buildGorillaHandler(rip.routes)
}

func buildGorillaHandler(routes []*Route) (http.Handler, error) {
	gorilla := mux.NewRouter()

	for _, route := range routes {
		gorilla.HandleFunc(route.Template(), createHandlerWrapper(route)).Methods(route.method)
	}

	return gorilla, nil
}
