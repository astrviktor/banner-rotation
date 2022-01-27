package internalhttp

import (
	"io"
	"net/http"
)

func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "OK\n")
}
