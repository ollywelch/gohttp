package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func NewUsersHandler(s Store) *UsersHandler {
	return &UsersHandler{
		store: s,
	}
}

func handleGetAuthenticatedUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(currentUserContextKey).(*User)
	if !ok || user == nil {
		writeJSON(w, "error getting current user", http.StatusInternalServerError)
		return
	}
	writeJSON(w, user, http.StatusOK)
}

func (uh *UsersHandler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := uh.store.GetUsers()
	if err != nil {
		writeJSON(w, "error getting users", http.StatusInternalServerError)
		return
	}
	writeJSON(w, users, http.StatusOK)
}

func (uh *UsersHandler) handleGetUsersById(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	idStr := params.ByName("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		writeJSON(w, fmt.Sprintf("error parsing id: %s", err.Error()), http.StatusBadRequest)
		return
	}

	user, err := uh.store.GetUserById(id)

	if err != nil {
		writeJSON(w, fmt.Sprintf("user with id=%d not found", id), http.StatusNotFound)
	}

	writeJSON(w, user, http.StatusOK)
}

func NewLoginHandler(s Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var c Credentials
		defer r.Body.Close()

		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			writeJSON(w, fmt.Sprintf("invalid request body: %s", err.Error()), http.StatusBadRequest)
			return
		}

		user, err := s.GetUserByName(c.Username)

		if err != nil {
			writeJSON(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		if user.Password != c.Password {
			writeJSON(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := NewJWT(user.Name)
		if err != nil {
			writeJSON(w, "authentication failed", http.StatusUnauthorized)
			return
		}
		writeJSON(w, token, http.StatusOK)
	}
}

func handleGetHealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, "healthy", http.StatusOK)
}
