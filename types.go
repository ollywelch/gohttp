package main

import "net/http"

type Store interface {
	GetUsers() ([]User, error)
	GetUserById(int) (*User, error)
	GetUserByName(string) (*User, error)
}

type InMemoryStore struct {
	users []User
}

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"-"`
}

type UsersHandler struct {
	store Store
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Middleware func(http.Handler) http.Handler
type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

type contextKey string

func (c contextKey) String() string {
	return "context key " + string(c)
}

type httpError struct {
	Status  int `json:"status"`
	Content any `json:"content"`
}
