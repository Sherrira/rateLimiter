package handler

import (
	"log"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received")
	w.Write([]byte("Hello, world!"))
}
