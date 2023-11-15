package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const authContextKey contextKey = "authKey"

type Middleware struct {
	Token *Token
}

// Middleware resolver
func GetMiddleware(ctx context.Context) *Middleware {
	out := &Middleware{Token: nil}

	token := ctx.Value(authContextKey)

	if token != nil {
		out.Token = token.(*Token)
	}

	return out
}

// Apply middleware
func ApplyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			next.ServeHTTP(w, r)
		} else {
			tokenRaw := strings.Split(authHeader, " ")[1]
			token := &Token{AuthId: tokenRaw}

			ctx := context.WithValue(r.Context(), authContextKey, token)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
