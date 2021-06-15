package rest

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Auth simulates an auth filter
func (r *REST) Auth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		if r.AuthKey != req.Header.Get("Authorization") {
			// Delegate request to the given handle

			//
			// Usually we would check if the user trying to do this has access to the requested resource
			//

			h(w, req, ps)
		} else {
			// Request Basic Authentication otherwise
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}
