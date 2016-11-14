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
	var result *rip.Request
	api := rip.NewRIP()
	api.Path("foo").Var("bar", "bardoc").GET().Do(func(api *rip.Route) {
		api.Handler(func(req *rip.Request, resp *rip.Response) {
			result = req
		}, "")
		api.Path("param").Param("baz", "bazdoc").Handler(func(req *rip.Request, resp *rip.Response) {
			result = req
		}, "")
		api.Var("baz", "bazdoc").Handler(func(req *rip.Request, resp *rip.Response) {
			result = req
		}, "")
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

		if err != nil {
			t.Errorf("failed request, %v", err)
		}

		if result.NumVars() != len(test.vars) {
			t.Errorf("expected number of vars to be %v but was %v", len(test.vars), result.NumVars())
			continue
		}

		for xVar, xVal := range test.vars {
			val, _ := result.GetVar(xVar)
			if val != xVal {
				t.Errorf("expected %v to be %v but was %v", xVar, xVal, val)
				continue
			}
		}
	}
}
