package server

import (
	"fmt"
	"io"
	"net/http"
)

func healthcheck(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "SanderServer healthy!")
}

// Start starts the server on the specified port
func Start(port string) {
	fmt.Println("SanderServer started on port " + port)

	http.HandleFunc("/", healthcheck)
	http.HandleFunc("/healthcheck", healthcheck)

	http.ListenAndServe(":8000", nil)
}
