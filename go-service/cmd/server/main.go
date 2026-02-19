// This is a simple Go HTTP server that listens on port 8080
package main

/*
This server has a single endpoint /health that responds with a message
indicating that the DocMind Go Service is running. It uses the net/http package
to handle HTTP requests and responses and prints a message to the console when
the server is running.
*/
import (
	"fmt"
	"net/http"
)

/*
Tha main function starts the HTTP server ans sets up the
/health endpoint to respond with a message indicating that the DocMind Go Service
is running.
It listens on port 8080 and prints a message to the console when the server is running.
*/
func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "DocMind Go Service is running")
	})

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
