package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	port := "8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	log.Printf("starting up on http://localhost:%s", port)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	// // Custom middleware for logging
	// r.Use(func(next http.Handler) http.Handler {
	// 	fn := func(w http.ResponseWriter, r *http.Request) {
	// 		log.Printf("got request: %+v\n", r)
	// 		next.ServeHTTP(w, r)
	// 	}
	// 	return http.HandlerFunc(fn)
	// })

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, map[string]string{"hello": "world"})
	})

	log.Fatal(http.ListenAndServe(":"+port, r))
}

func JSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	b, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		JSON(w, map[string]string{"error": err.Error()})
		return
	}
	w.Write(b)
}
