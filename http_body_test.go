package rip_test

import (
	"bytes"
	"fmt"
	"github.com/creichlin/rip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPJsonData(t *testing.T) {
	result := ""
	api := rip.NewRIP()
	api.Path("foo").GET().Target(map[string]interface{}{}).Handler(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		data := rip.Body(req)
		result = fmt.Sprintf("%#v", data)
	}), "")
	api.Path("foo").POST().Target(map[string]interface{}{}).Handler(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		data := rip.Body(req)
		result = fmt.Sprintf("%#v", data)
	}), "")
	api.Path("foo").PUT().Target(map[string]interface{}{}).Handler(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		data := rip.Body(req)
		result = fmt.Sprintf("%#v", data)
	}), "")
	handler, err := api.RootHandler()
	if err != nil {
		t.Errorf("erroneous route %v", err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	url := server.URL

	tests := []struct {
		method      string
		contentType string
		body        string
		result      interface{}
	}{
		{"GET", "", "", `""`},
		{"GET", "application/json", `{"foo": "bar"}`, `&map[string]interface {}{"foo":"bar"}`},
		{"POST", "application/json", `{"foo": "bar"}`, `&map[string]interface {}{"foo":"bar"}`},
		{"PUT", "application/json", `{"foo": "bar"}`, `&map[string]interface {}{"foo":"bar"}`},
		{"GET", "application/vnd.x.y+json", `{"foo": "bar"}`, `&map[string]interface {}{"foo":"bar"}`},
		{"GET", "application/json", `{"föö": "bär"}`, `&map[string]interface {}{"föö":"bär"}`},
		{"GET", "application/json; charset=utf-8", `{"föö": "bär"}`, `&map[string]interface {}{"föö":"bär"}`},
		{"GET", "application/json; charset=iso-8859-1", `{"föö": "bär"}`, `&map[string]interface {}{"föö":"bär"}`},
	}

	for _, test := range tests {
		var body io.Reader
		if test.body != "" {
			body = bytes.NewBufferString(test.body)
		}

		request, _ := http.NewRequest(test.method, url+"/foo", body)
		request.Header.Set("Content-Type", test.contentType)
		result = ""
		_, _ = http.DefaultClient.Do(request)
		if result != test.result {
			t.Errorf("expected data to be %v but was %v\n%v", test.result, result, test)
		}
	}
}
