package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
)

const (
	port                  = ":3000"
	currentUserContextKey = contextKey("auth.currentUser")
)

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func main() {
	router := httprouter.New()

	store := NewInMemoryStore()

	authMiddleware := NewAuthMiddleware(store)
	usersHandler := NewUsersHandler(store)

	router.HandlerFunc("POST", "/auth/login", usersHandler.handleLogin)
	router.HandlerFunc("GET", "/auth/me", AdaptHandlerFunc(handleGetAuthenticatedUser, authMiddleware))
	router.HandlerFunc("GET", "/users", AdaptHandlerFunc(usersHandler.handleGetUsers, authMiddleware))
	router.HandlerFunc("GET", "/users/:id", AdaptHandlerFunc(usersHandler.handleGetUsersById, authMiddleware))
	router.HandlerFunc("GET", "/healthz", handleGetHealthCheck)

	wrappedRouter := AdaptHandler(router, LogRequests)

	slog.Info(fmt.Sprintf("http server listening at %s", port))
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
		slog.Info(fmt.Sprintf("%s %s %v", r.Method, r.URL.Path, time.Since(start)))
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
