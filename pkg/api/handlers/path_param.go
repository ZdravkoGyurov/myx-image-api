package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func pathParam(request *http.Request, param string) string {
	value, ok := mux.Vars(request)[param]
	if !ok {
		panic("could not find path param")
	}
	return value
}
