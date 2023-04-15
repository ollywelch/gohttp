package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type UsersHandler struct {
	store Store
}

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

func handleGetToken(w http.ResponseWriter, r *http.Request) {
	token, err := NewJWT("Olly")
	if err != nil {
		writeJSON(w, fmt.Errorf("authentication failed"), http.StatusUnauthorized)
		return
	}
	writeJSON(w, token, http.StatusOK)
}

func handleGetHealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, "healthy", http.StatusOK)
}
