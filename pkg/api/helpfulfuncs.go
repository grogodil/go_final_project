package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

func afterNow(date, now time.Time) bool {
	if date.Format(DateFormat) > now.Format(DateFormat) {
		return true
	} else {
		return false
	}
}
