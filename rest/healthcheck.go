package rest

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Healthcheck - checks if the service is up
func (r *REST) Healthcheck(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}
