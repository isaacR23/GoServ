package main

import (
	"fmt"
	"net/http"
)

type apiHandler struct{}

func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func main() {
	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)

	mux.Handle("/api/", apiHandler{})
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {

		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}

	// fmt.Println("Server is running on http://localhost:8080")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNotFound) // change for http.StatusOK for code 200 OK
			return
		}
		next.ServeHTTP(w, r)
	})
}
