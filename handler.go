package rip

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type RIPResponse struct {
	StatusCode int
	Data       interface{}
}

type Variables map[string]string

func (v Variables) GetVar(name string) (string, error) {
	val, exists := v[name]
	if !exists {
		return "", fmt.Errorf("Variable %v not declared in route", name)
	}
	return val, nil
}

func (v Variables) MustGetVar(name string) string {
	val, err := v.GetVar(name)
	if err != nil {
		panic(err)
	}
	return val
}

func (v Variables) NumVars() int {
	return len(v)
}

func createParseVarsHandler(nextHandler http.Handler) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		variables := Variables{}

		for k, v := range mux.Vars(request) {
			variables[k] = v
		}

		route := request.Context().Value("rip-route").(*Route)

		for _, v := range route.queryParameters {
			variables[v.name] = request.URL.Query().Get(v.name)
		}

		nextHandler.ServeHTTP(response,
			request.WithContext(context.WithValue(request.Context(), "rip-variables", variables)))
	}
}

// if content type is application/json and request body is valid json
// parse it and store it in a new instance of targetType that got
// defined in the route and store it in the context
func createParseBodyHandler(nextHandler http.Handler) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		contentType := request.Header.Get("Content-Type")

		route := request.Context().Value("rip-route").(*Route)

		data, err := parseData(contentType, route.target, request.Body)
		if err != nil {
			log.Printf("parse vars handler %v", err)
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte(fmt.Sprintf("Error parsing request body, %v", err))) // nolint: errcheck
			return
		}
		nextHandler.ServeHTTP(response,
			request.WithContext(context.WithValue(request.Context(), "rip-body", data)))
	}
}

func createResponseWriter(nextHandler http.Handler) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		ripResponse := &RIPResponse{
			StatusCode: http.StatusOK,
		}
		nextHandler.ServeHTTP(response,
			request.WithContext(context.WithValue(request.Context(), "rip-response", ripResponse)))

		// force all answers to be json for now
		response.Header().Set("Content-Type", "application/json; charset=UTF-8")

		jsonData, err := json.Marshal(ripResponse.Data)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"response": "internal server error"}`)) // nolint: errcheck
			log.Printf("failed to marshal response %v, %v", ripResponse.Data, err)
			return
		}
		response.WriteHeader(ripResponse.StatusCode)
		_, err = response.Write(jsonData)
		if err != nil {
			route := request.Context().Value("rip-route").(*Route)
			log.Printf("failed to write response for request %v, %v", route.Template(), err)
		}
	}
}

func createHandlerWrapper(route *Route) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		createResponseWriter(
			createParseBodyHandler(
				createParseVarsHandler(
					route.handler,
				),
			),
		)(w, r.WithContext(context.WithValue(r.Context(), "rip-route", route)))
	}
}

func parseData(contentType string, targetType interface{}, data io.Reader) (interface{}, error) {
	ct := removeMimeVendor(contentType)
	ct = removeMimeEncoding(ct)

	if ct == "application/json" {
		if targetType == nil {

			return nil, fmt.Errorf("Request does not accept a body, %v", targetType)
		}

		targetData := reflect.New(reflect.TypeOf(targetType)).Interface()

		decoder := json.NewDecoder(data)
		err := decoder.Decode(&targetData)

		if err != nil {
			return nil, err
		}
		return targetData, nil
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

func DocHandler(rip *rip) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
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
		Response(request).Data = resp
	})
}
