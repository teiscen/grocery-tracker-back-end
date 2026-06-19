//go:build ignore

package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handler) ScanReceipt(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "image field missing", http.StatusBadRequest)
		return
	}
	defer file.Close()

	result, err := h.ScanService.Process(file)
	if err != nil {
		// log the real error server side, don't expose it to the client
		log.Printf("scan error: %v", err)
		http.Error(w, "scan failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) RegisterScanRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/scan", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.ScanReceipt(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}