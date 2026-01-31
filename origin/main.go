package main

import (
	"log"
	"net/http"
)

func main() {
	// Serves everything in the /storage folder
	fs := http.FileServer(http.Dir("./storage"))
	http.Handle("/", fs)

	log.Println("Origin Server started on :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
