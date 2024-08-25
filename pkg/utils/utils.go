package utils

import (
	"log"
	"net/http"
)

func HandleError(w http.ResponseWriter, message string, code int) {
	log.Printf("Error: %s", message)
	http.Error(w, message, code)
}
