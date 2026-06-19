package handlers

import "net/http"

// Registerable must be implemented by all resource handlers
type Registerable interface {
	Register(mux *http.ServeMux)
}

