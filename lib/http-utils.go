package lib

import (
	"encoding/json"
	"net/http"
)

// RenderJSON write data as a json
func RenderJSON(w http.ResponseWriter, data interface{}, err error) {
	if isInvalid(w, err) {
		return
	}
	response, err := json.Marshal(data)
	if isInvalid(w, err) {
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(response)
}

func isInvalid(w http.ResponseWriter, err error) bool {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}

func header(r *http.Request, key string) (string, bool) {
	if r.Header == nil {
		return "", false
	}
	if candidate := r.Header[key]; len(candidate) > 0 {
		return candidate[0], true
	}
	return "", false
}
