package utils

import "net/http"

func RespondeWithError(w http.ResponseWriter, code int, message string) {
	RespondeWithJSON(w, code, map[string]any{"error": message})
}
