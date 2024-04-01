package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileserverHits++
        next.ServeHTTP(w, r)
    })
}



func main() {
	const port = "8080"

	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	var apiCfg apiConfig
	handler := http.FileServer(http.Dir("./assets/"))
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(handler))

	mux.Handle("/app", apiCfg.middlewareMetricsInc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to Chirpy"))
	})))


	mux.HandleFunc("GET /metrics",func(curl http.ResponseWriter, req *http.Request){
		// note, you get information from the client aka the response when you send something
		// with the req. 
		curl.WriteHeader(http.StatusOK)
		var hits string = fmt.Sprintf("Hits: %v", apiCfg.fileserverHits)
		curl.Write([]byte(hits))
	})

	mux.HandleFunc("/reset",func(curl http.ResponseWriter, req *http.Request){
		// note, you get information from the client aka the response when you send something
		// with the req. 
		apiCfg.fileserverHits = 0
		curl.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("GET /healthz",func(curl http.ResponseWriter, req *http.Request){

		curl.Header().Set("Content-Type", "text/plain; charset=utf-8")
		curl.WriteHeader(http.StatusOK)
		curl.Write([]byte("OK"))
	})

	// Serve the assets directory under the path /app/assets/
	logoDir := http.FileServer(http.Dir("./assets/"))
	mux.Handle("/app/assets/", http.StripPrefix("/app/assets/", logoDir))

	log.Printf("Serving on port: http://localhost:%s\n", port)
	log.Fatal(srv.ListenAndServe()) // This exe the serv. Anything after this line is never ran
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