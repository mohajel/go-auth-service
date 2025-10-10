package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type key int

const userIDKey key = 0

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		userID := parts[1]

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey)
	if userID == nil {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Hello user %s", userID)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/protected", Middleware(http.HandlerFunc(ProtectedHandler)))

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", mux)
}
