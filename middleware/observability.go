package middleware

import (
	"log"
	"net/http"
)

func ObservabilityMiddleware(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("starting request: %s - %s", r.Method, r.URL)
		next.ServeHTTP(w, r)
		log.Printf("finished request: %s - %s", r.Method, r.URL)
	})
}
