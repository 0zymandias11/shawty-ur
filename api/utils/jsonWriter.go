package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	response := map[string]interface{}{
				"error": msg,
	}

	// If encoding fails, you want to *log* it but avoid recursive failures.
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to write JSON response", "err", err)
	}
}
