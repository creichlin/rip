package rip_test

import (
	"github.com/creichlin/rip"
	"net/http"
	"net/http/httptest"
	"testing"
)

type vars map[string]string

type varTest struct {
	url  string
	vars map[string]string
}

func TestHTTPVarData(t *testing.T) {
	var result *http.Request
	reqHandler := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		result = req
	})
	api := rip.NewRIP()
	api.Path("foo").Var("bar", "bardoc").GET().Do(func(api *rip.Route) {
		api.Handler(reqHandler, "")
		api.Path("param").Param("baz", "bazdoc").Handler(reqHandler, "")
		api.Var("baz", "bazdoc").Handler(reqHandler, "")
	})

	tests := []varTest{
		{"/foo/aaa", vars{"bar": "aaa"}},
		{"/foo/A  B", vars{"bar": "A  B"}},
		{"/foo/aaa/1234", vars{"bar": "aaa", "baz": "1234"}},
		{"/foo/aaa/param", vars{"bar": "aaa", "baz": ""}},
		{"/foo/aaa/param?xxx=yyy", vars{"bar": "aaa", "baz": ""}},
		{"/foo/aaa/param?baz=yyy", vars{"bar": "aaa", "baz": "yyy"}},
	}

	handler, err := api.RootHandler()
	if err != nil {
		t.Errorf("erroneous route %v", err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	url := server.URL

	for _, test := range tests {
		req, _ := http.NewRequest("GET", url+test.url, nil)
		_, err = http.DefaultClient.Do(req)
		variables := rip.Vars(result)

		if err != nil {
			t.Errorf("failed request, %v", err)
		}

		if variables.NumVars() != len(test.vars) {
			t.Errorf("expected number of vars to be %v but was %v", len(test.vars), variables.NumVars())
			continue
		}

		for xVar, xVal := range test.vars {
			val, _ := variables.GetVar(xVar)
			if val != xVal {
				t.Errorf("expected %v to be %v but was %v", xVar, xVal, val)
				continue
			}
		}
	}
}
