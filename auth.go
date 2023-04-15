package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("supersecretpassword")

func NewAuthMiddleware(s Store) MiddlewareFunc {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			currentUser, err := authenticateRequest(r, s)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), currentUserContextKey, currentUser)

			rWithCtx := r.WithContext(ctx)

			hf.ServeHTTP(w, rWithCtx)
		}
	}
}

func NewJWT(name string) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject: name,
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

func authenticateRequest(r *http.Request, s Store) (*User, error) {
	token, err := getTokenFromHeader(r)

	if err != nil {
		return nil, err
	}

	claims, err := ValidateJWT(token)

	if err != nil {
		return nil, err
	}

	username, err := claims.GetSubject()

	if err != nil {
		return nil, err
	}

	return s.GetUserByName(username)
}

func getTokenFromHeader(r *http.Request) (string, error) {
	headerVal := r.Header.Get("Authorization")

	if headerVal == "" {
		return "", fmt.Errorf("not authenticated")
	}

	splitHeader := strings.Split(headerVal, "Bearer ")
	if len(splitHeader) < 2 {
		return "", fmt.Errorf("invalid authorization header")
	}
	return splitHeader[1], nil
}
