package handlers

import (
	"net/http"
	"strings"

	"receiptTracker/services"
)

type LocationHandler struct {
	Service *services.LocationServices
}

func (h *LocationHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/location", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetLocations(w, r)
		case http.MethodPost:
			h.CreateLocation(w, r)
		case http.MethodDelete:
			h.DeleteLocations(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/location/", func(w http.ResponseWriter, r *http.Request) {
		// handle nested route
		if strings.HasSuffix(r.URL.Path, "/products") {
			switch r.Method {
			case http.MethodGet:
				h.GetProductsByLocation(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		}

		switch r.Method {
		case http.MethodGet:
			h.GetLocation(w, r)
		case http.MethodDelete:
			h.DeleteLocation(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (h *LocationHandler) GetProductsByLocation(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/products")
	id, err := parseID(path, "/locations/")
	if err != nil {
		writeError(w, http.StatusBadRequest,
			"invalid location id", err,
		)
	}
	locations, err := h.Service.GetProductsByLocation(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError,
			"failed to get products", err)
		return
	}

	writeJson(w, http.StatusOK, locations)
}

func (h *LocationHandler) DeleteLocations(w http.ResponseWriter, r *http.Request) {
	err := h.Service.DeleteLocations()
	if err != nil {
		writeError(w, http.StatusInternalServerError,
			"delete location error", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *LocationHandler) DeleteLocation(w http.ResponseWriter, r *http.Request) {
	// extract id from URL e.g. /locations/1
	id, err := parseID(r.URL.Path, "/location/")
	if err != nil {
		writeError(w, http.StatusBadRequest,
			"invalid location id", err)
		return
	}

	err = h.Service.DeleteLocation(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError,
			"failed to delete location", err)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204 -- success, nothing to return
}

func (h *LocationHandler) GetLocations(w http.ResponseWriter, r *http.Request) {
	locations, err := h.Service.GetLocations()
	if err != nil {
		writeError(w, http.StatusInternalServerError,
			"get locations error", err)
		return
	}
	writeJson(w, http.StatusOK, locations)
}

func (h *LocationHandler) GetLocation(w http.ResponseWriter, r *http.Request) {
	// extract id from url (/locations/1)
	id, err := parseID(r.URL.Path, "/location/")
	if err != nil {
		writeError(w, http.StatusBadRequest,
			"invalid location id", err)
		return
	}

	location, err := h.Service.GetLocation(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError,
			"failed to get location", err)
		return
	}

	writeJson(w, http.StatusOK, location)
}

func (h *LocationHandler) CreateLocation(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeError(w, http.StatusBadRequest,
			"invalid request body", err)
		return
	}
	location, err := h.Service.CreateLocation(body.Name, body.Type)
	if err != nil {
		writeError(w, http.StatusInternalServerError,
			"failed to create location", err)
		return
	}
	writeJson(w, http.StatusCreated, location)
}
