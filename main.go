package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const port = ":3000"

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Logger struct {
	h http.Handler
}

type httpError struct {
	Status  int `json:"status"`
	Content any `json:"content"`
}

type handler func(http.ResponseWriter, *http.Request) error

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
	mux := http.NewServeMux()
	mux.HandleFunc("/users", makeHttpHandler(handleGetUsers))

	loggedMux := NewLogger(mux)

	log.Fatal(http.ListenAndServe(port, loggedMux))
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, users, http.StatusOK)
}

func makeHttpHandler(h handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.h.ServeHTTP(w, r)
	log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
}

func NewLogger(h http.Handler) *Logger {
	return &Logger{h}
}

func writeJSON(w http.ResponseWriter, v any, status int) error {
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
		return err
	}

	return nil
}
