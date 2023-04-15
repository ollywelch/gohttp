package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

const (
	port                  = ":3000"
	currentUserContextKey = contextKey("auth.currentUser")
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type httpError struct {
	Status  int `json:"status"`
	Content any `json:"content"`
}

type Middleware func(http.Handler) http.Handler
type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

type contextKey string

func (c contextKey) String() string {
	return "context key " + string(c)
}

func main() {
	router := httprouter.New()

	store := NewInMemoryStore()

	ensureAuth := NewAuthMiddleware(store)

	usersHandler := NewUsersHandler(store)
	router.HandlerFunc("POST", "/auth/login", handleGetToken)
	router.HandlerFunc("GET", "/users", AdaptHandlerFunc(usersHandler.handleGetUsers, ensureAuth))
	router.HandlerFunc("GET", "/users/:id", AdaptHandlerFunc(usersHandler.handleGetUsersById, ensureAuth))
	router.HandlerFunc("GET", "/auth/me", AdaptHandlerFunc(handleGetAuthenticatedUser, ensureAuth))
	router.HandlerFunc("GET", "/healthz", handleGetHealthCheck)

	wrappedRouter := AdaptHandler(router, LogRequests)

	log.Println("http server listening at", port)
	log.Fatal(http.ListenAndServe(port, wrappedRouter))
}

func AdaptHandler(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

func AdaptHandlerFunc(hf http.HandlerFunc, mfs ...MiddlewareFunc) http.HandlerFunc {
	for _, mf := range mfs {
		hf = mf(hf)
	}
	return hf
}

func LogRequests(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func writeJSON(w http.ResponseWriter, v any, status int) {
	var res any

	if status != http.StatusOK {
		res = httpError{
			Status:  status,
			Content: v,
		}
		w.WriteHeader(status)
	} else {
		res = v
	}

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
