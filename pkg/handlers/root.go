package handlers

import (
	"log"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("OK.")
	w.WriteHeader(http.StatusOK)
}
