package server

import (
	"log"
	"net/http"
)

func serverError(w http.ResponseWriter, err error) {
	log.Println(err)
	errorResponse(w, http.StatusInternalServerError, "internal error")
}

func errorResponse(w http.ResponseWriter, code int, errors interface{}) {
	writerJSON(w, code, M{"errors": errors})
}
