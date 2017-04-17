package rip_test

import (
	"github.com/creichlin/rip"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPUnconfiguredVarData(t *testing.T) {
	varTryHandler := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		vars := req.Context().Value("rip-variables").(rip.Variables)
		for _, varr := range []string{"123", "boo", "  ", "-23_"} {
			_, err := vars.GetVar(varr)
			if err == nil {
				t.Errorf("Could read undefined var %v", varr)
			}
		}
	})

	api := rip.NewRIP()
	api.Path("foo").Var("bar", "bardoc").GET().Do(func(api *rip.Route) {
		api.Handler(varTryHandler, "")
		api.Path("param").Param("baz", "bazdoc").Handler(varTryHandler, "")
		api.Var("baz", "bazdoc").Handler(varTryHandler, "")
	})

	tests := []string{
		"/foo/aaa",
		"/foo/A  B",
		"/foo/aaa/1234",
		"/foo/aaa/param",
		"/foo/aaa/param?xxx=yyy",
		"/foo/aaa/param?baz=yyy",
	}

	handler, err := api.RootHandler()
	if err != nil {
		t.Errorf("erroneous route %v", err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	url := server.URL

	for _, test := range tests {
		req, _ := http.NewRequest("GET", url+test, nil)
		_, _ = http.DefaultClient.Do(req)
	}
}
