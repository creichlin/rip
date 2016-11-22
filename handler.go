package rip

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strings"
)

type Request struct {
	variables   map[string]string
	Data        interface{}
	HttpRequest *http.Request
	route       *Route
}

type Response struct {
	StatusCode   int
	Data         interface{}
	HttpResponse http.ResponseWriter
}

func (r *Request) GetVar(name string) (string, error) {
	val, exists := r.variables[name]
	if !exists {
		return "", fmt.Errorf("Variable %v not declared in route", name)
	}
	return val, nil
}

func (r *Request) MustGetVar(name string) string {
	val, err := r.GetVar(name)
	if err != nil {
		panic(err)
	}
	return val
}

func (r *Request) NumVars() int {
	return len(r.variables)
}

type Handler func(request *Request, response *Response)

func createParseVarsHandler(nextHandler Handler) Handler {
	return func(request *Request, response *Response) {
		request.variables = map[string]string{}

		for k, v := range mux.Vars(request.HttpRequest) {
			request.variables[k] = v
		}

		for _, v := range request.route.queryParameters {
			request.variables[v.name] = request.HttpRequest.URL.Query().Get(v.name)
		}
		nextHandler(request, response)
	}
}

func createParseBodyHandler(nextHandler Handler) Handler {
	return func(request *Request, response *Response) {
		contentType := request.HttpRequest.Header.Get("Content-Type")
		data, err := parseData(contentType, request.HttpRequest.Body)
		if err != nil {
			fmt.Printf("parse vars handler %v\n", err)
			response.StatusCode = http.StatusBadRequest
			response.Data = map[string]string{"error": err.Error()}
			return
		}
		request.Data = data
		nextHandler(request, response)
	}
}

func createResponseWriter(nextHandler Handler) Handler {
	return func(request *Request, response *Response) {
		nextHandler(request, response)

		// force all answers to be json for now
		response.HttpResponse.Header().Set("Content-Type", "application/json; charset=UTF-8")

		jsonData, err := json.Marshal(response.Data)
		if err != nil {
			response.HttpResponse.WriteHeader(http.StatusInternalServerError)
			response.HttpResponse.Write([]byte(`{"response": "internal server error"}`))
			log.Printf("failed to marshal response %v, %v", response.Data, err)
			return
		}
		response.HttpResponse.WriteHeader(response.StatusCode)
		_, err = response.HttpResponse.Write(jsonData)
		if err != nil {
			log.Printf("failed to write response for request %v, %v", request.route.Template(), err)
		}
	}
}

func createHandlerWrapper(route *Route) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		request := &Request{
			HttpRequest: r,
			route:       route,
		}

		response := &Response{
			HttpResponse: w,
			StatusCode:   http.StatusOK,
		}

		createResponseWriter(
			createParseBodyHandler(
				createParseVarsHandler(
					route.handler,
				),
			),
		)(request, response)
	}
}

func parseData(contentType string, data io.Reader) (interface{}, error) {
	ct := removeMimeVendor(contentType)
	ct = removeMimeEncoding(ct)

	if ct == "application/json" {
		decoder := json.NewDecoder(data)
		mapData := map[string]interface{}{}

		err := decoder.Decode(&mapData)
		if err != nil {
			return nil, err
		}
		return mapData, nil
	} else if ct == "" {
		return "", nil
	}

	return nil, fmt.Errorf("Invalid content type: %v", contentType)
}

func removeMimeEncoding(ct string) string {
	index := strings.Index(ct, ";")
	if index != -1 {
		return ct[:index]
	}
	return ct
}

func removeMimeVendor(ct string) string {
	slashIndex := strings.Index(ct, "/")
	plusIndex := strings.Index(ct, "+")
	if slashIndex != -1 && plusIndex > slashIndex {
		return ct[:slashIndex+1] + ct[plusIndex+1:]
	}
	return ct
}

func DocHandler(rip *rip) Handler {
	return func(request *Request, response *Response) {
		resp := map[string]interface{}{}
		links := []map[string]interface{}{}

		for _, route := range rip.Routes() {
			li := map[string]interface{}{}
			li["href"] = route.Template() + route.QueryTemplate()
			li["method"] = route.Method()
			li["description"] = route.doc

			parameters := []interface{}{}

			for _, v := range route.parameters {
				vars := map[string]string{
					"name":        v.Name(),
					"description": v.Doc(),
				}
				parameters = append(parameters, vars)
			}
			li["parameters"] = parameters
			links = append(links, li)
		}

		resp["links"] = links
		response.Data = resp
	}
}
