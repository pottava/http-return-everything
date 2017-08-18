package lib

import (
	"encoding/json"
	"net/http"
	"strings"
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

func isInvalid(w http.ResponseWriter, err error) (invalid bool) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return
}

// Header retrive HTTP header values with specified key
func Header(r *http.Request, key string) (values []string, found bool) {
	if r.Header == nil {
		return
	}
	for k, v := range r.Header {
		if strings.EqualFold(k, key) && len(v) > 0 {
			return v, true
		}
	}
	return
}
