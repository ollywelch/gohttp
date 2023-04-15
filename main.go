package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

var users = []User{
	{
		Id:   1,
		Name: "Olly",
	},
	{
		Id:   2,
		Name: "Jeff",
	},
}

func main() {
	router := httprouter.New()

	ensureAuth := EnsureAuthentication()

	router.HandlerFunc("GET", "/users", AdaptHandlerFunc(handleGetUsers, ensureAuth))
	router.HandlerFunc("GET", "/users/:id", AdaptHandlerFunc(handleGetUsersById, ensureAuth))
	router.HandlerFunc("GET", "/auth/me", AdaptHandlerFunc(handleGetCurrentUser, ensureAuth))

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

func AdaptHandlerFunc(h http.HandlerFunc, middlewares ...MiddlewareFunc) http.HandlerFunc {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

func LogRequests(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func EnsureAuthentication() MiddlewareFunc {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			currentUser, err := authenticateRequest(r)
			if err != nil {
				http.Error(w, "not authenticated", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), currentUserContextKey, currentUser)

			rWithCtx := r.WithContext(ctx)

			hf.ServeHTTP(w, rWithCtx)
		}
	}
}

func authenticateRequest(r *http.Request) (*User, error) {
	//TODO: Some auth
	return &users[0], nil
}

func handleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(currentUserContextKey).(*User)
	if !ok || user == nil {
		writeJSON(w, "error getting current user", http.StatusInternalServerError)
		return
	}
	writeJSON(w, user, http.StatusOK)
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, users, http.StatusOK)
}

func handleGetUsersById(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	idStr := params.ByName("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		writeJSON(w, fmt.Sprintf("error parsing id: %s", err.Error()), http.StatusBadRequest)
		return
	}

	for _, user := range users {
		if user.Id == id {
			writeJSON(w, user, http.StatusOK)
			return
		}
	}

	writeJSON(w, fmt.Sprintf("user with id=%d not found", id), http.StatusNotFound)
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
