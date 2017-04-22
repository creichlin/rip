package rip

import "net/http"

func Vars(r *http.Request) Variables {
	vars := r.Context().Value("rip-variables")
	tVars, ok := vars.(Variables)
	if !ok {
		panic("rip-variables must be of type Variables. Don't overwritte.")
	}
	return tVars
}

func Body(r *http.Request) interface{} {
	return r.Context().Value("rip-body")
}

func Response(r *http.Request) *RIPResponse {
	resp := r.Context().Value("rip-response")
	tResp, ok := resp.(*RIPResponse)
	if !ok {
		panic("rip-response must be of type *RIPResponse. Don't overwritte.")
	}
	return tResp
}
