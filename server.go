package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/calvine/urshrt/middleware"
)

const embiggenEndpoint = "/e/"

type server struct {
	service URLShortenerService
}

func NewServer(service URLShortenerService) *server {
	return &server{
		service: service,
	}
}

func (server) rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(400)
		return
	}
	w.Write([]byte("Bok"))
}

func (s *server) shrtHandler(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only allow HTTP POST
		if r.Method != http.MethodPost {
			logger.Printf("got wrong http method expected POST: %s", r.Method)
			w.WriteHeader(400)
			return
		}
		url := r.FormValue("url")
		logger.Printf("got url from form data: %s", url)
		// r.URL.Query()
		// queryParams, err := url.ParseQuery(r.URL.RawQuery)
		// if err != nil {
		// 	logger.Printf("failed to parse query on url: %s - %v", r.URL, err)
		// 	w.WriteHeader(400)
		// 	return
		// }
		// u := queryParams.Get("u")
		// if len(u) == 0 {
		// 	logger.Printf("no url to shorten provided: %s", r.URL)
		// 	w.WriteHeader(400)
		// 	return
		// }
		// uu, err := url.QueryUnescape(u)
		// if err != nil {
		// 	logger.Printf("failed to unescape url from query param: %s - %v", uu, err)
		// 	w.WriteHeader(400)
		// 	return
		// }
		// logger.Printf("got url: %s", uu)
		shrtKey, err := s.service.Shorten(r.Context(), logger, url, 10)
		if err != nil {
			logger.Printf("failed to generate short key for URL %s - %v", url, err)
			w.WriteHeader(500)
			return
		}
		// TODO: return json and set content type to application/json
		w.Write([]byte(shrtKey))
	}
}

func (s *server) embiggenHandler(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			logger.Printf("got wrong http method expected GET: %s", r.Method)
			w.WriteHeader(400)
			return
		}
		logger.Printf("got call to embiggen endpoint: %s", r.URL.Path)
		// so to try and be clever  this endpoint should always be the same with the /e/ in the begenning...
		// so we can just shop off the first 3 characters and use the rest as the key, with some kind of limit to be safe...
		// like 25 characters
		key, found := strings.CutPrefix(r.URL.Path, embiggenEndpoint)
		if !found {
			w.WriteHeader(400)
			logger.Printf("something fishy with the url path... %s", r.URL.Path)
			return
		}
		principal, err := s.service.Embiggen(r.Context(), logger, key)
		if err != nil {
			logger.Printf("failed to embiggen the key %s - %v", key, err)
			w.WriteHeader(500)
			return
		}
		logger.Printf("redirecting to %s", principal)
		http.Redirect(w, r, principal, http.StatusTemporaryRedirect)
	}
}

func (s *server) InitServer(logger *log.Logger, addr string) error {
	mux := http.NewServeMux()
	mux.Handle(embiggenEndpoint, middleware.ObservabilityMiddleware(logger, s.embiggenHandler(logger)))
	mux.Handle("/shrt", middleware.ObservabilityMiddleware(logger, s.shrtHandler(logger)))
	rh := http.HandlerFunc(s.rootHandler)
	mux.Handle("/", middleware.ObservabilityMiddleware(logger, rh))
	log.Printf("starting server: %s", addr)
	return http.ListenAndServe(addr, mux)
}
