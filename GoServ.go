package main

import (
	"log"
	"net/http"
)



func main() {
	const port = "8080"

	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	// Serve the assets directory under the path /app/assets/
	logoDir := http.FileServer(http.Dir("./assets/"))
	mux.Handle("/app/assets/", http.StripPrefix("/app/assets/", logoDir))
	// Serve index.html at the path /app
	mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./index.html")
	})

	mux.HandleFunc("/healthz",func(curl http.ResponseWriter, req *http.Request){

		curl.Header().Set("Content-Type", "text/plain; charset=utf-8")
		curl.WriteHeader(http.StatusOK)
		curl.Write([]byte("OK"))
	})

	log.Printf("Serving on port: http://localhost:%s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}