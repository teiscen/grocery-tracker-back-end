package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"strconv"
)

func decodeBody(r *http.Request, dst interface{}) error {
	return json.NewDecoder(r.Body).Decode(dst)
}

func writeJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string, err error){
    log.Printf("%s: %v", msg, err)
    http.Error(w, msg, status)
}

func parseID(path string, prefix string) (int, error) {
    idStr := strings.TrimPrefix(path, prefix)
    return strconv.Atoi(idStr)
}