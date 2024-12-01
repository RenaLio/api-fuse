package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/v1/completions", handleCompletion)
	http.HandleFunc("/v1/chat/completions", HandleChat)
	http.HandleFunc("/v1/models", handleModel)

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
