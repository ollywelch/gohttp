package main

import (
	"context"
	"net/http"
)

func NewAuthMiddleware(s Store) MiddlewareFunc {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			currentUser, err := authenticateRequest(r, s)
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

func authenticateRequest(r *http.Request, s Store) (*User, error) {
	return s.GetUserByName("Olly")
}
