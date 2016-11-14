package rip

import (
	"fmt"
	"testing"
)

var dummyHandler = func(q *Request, n *Response) {}

func TestInspectRoutes(t *testing.T) {
	rip := NewRIP()
	rip.Path("foo").Do(func(rip *Route) {
		rip.GET().Handler(dummyHandler, "")
		rip.POST().Handler(dummyHandler, "")
		rip.Path("bar").Do(func(rip *Route) {
			rip.GET().Handler(dummyHandler, "")
			rip.PUT().Handler(dummyHandler, "")
			rip.DELETE().Handler(dummyHandler, "")
		})
	})

	result := ""

	for _, route := range rip.Routes() {
		result += fmt.Sprintf("%v %v\n", route.Template(), route.Method())
	}

	if result != `/foo GET
/foo POST
/foo/bar GET
/foo/bar PUT
/foo/bar DELETE
` {
		t.Errorf("Wrong routes calculated")
	}
}

func TestInspectVar(t *testing.T) {
	rip := NewRIP()
	rip.Path("foo").Var("bar", "id of foo item").GET().Handler(dummyHandler, "")

	template := rip.Routes()[0].Template()
	if template != "/foo/{bar}" {
		t.Errorf("wrong template %v", template)
	}

	if len(rip.Routes()[0].parameters) != 1 {
		t.Errorf("route must have 1 variable")
	}

	v := rip.Routes()[0].parameters[0]

	if v.Name() != "bar" {
		t.Errorf("expected variable to be 'bar' but is '%v'", v.Name())
	}

	if v.Doc() != "id of foo item" {
		t.Errorf("expected variable to have doc 'id of foo item' but is '%v'", v.Doc())
	}
}

func TestInspectParam(t *testing.T) {
	rip := NewRIP()
	rip.Path("foo").Param("bar", "id of foo item").GET().Handler(dummyHandler, "")

	template := rip.Routes()[0].Template()
	if template != "/foo" {
		t.Errorf("wrong template %v", template)
	}

	if len(rip.Routes()[0].parameters) != 1 {
		t.Errorf("route must have 1 variable")
	}

	v := rip.Routes()[0].parameters[0]

	if v.Name() != "bar" {
		t.Errorf("expected variable to be 'bar' but is '%v'", v.Name())
	}

	if v.Doc() != "id of foo item" {
		t.Errorf("expected variable to have doc 'id of foo item' but is '%v'", v.Doc())
	}
}
