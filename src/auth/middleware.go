package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/glimpzio/backend/misc"
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

type Auth0Config struct {
	Auth0Domain       string
	Auth0ClientId     string
	Auth0ClientSecret string
}

// Verify a token
func VerifyToken(ctx context.Context, token string, config *Auth0Config) (*Token, error) {
	provider, err := oidc.NewProvider(ctx, fmt.Sprintf("https://%s/", config.Auth0Domain))
	if err != nil {
		return nil, err
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: config.Auth0ClientId})

	idToken, err := verifier.Verify(ctx, token)
	if err != nil {
		return nil, err
	}

	var claims struct {
		Email string `json:"email"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return nil, err
	}

	return &Token{AuthId: idToken.Subject, Email: claims.Email}, nil
}

// Apply middleware
func ApplyMiddleware(logger *misc.Logger, next http.Handler, config *Auth0Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gc, _ := misc.GinContextFromContext(r.Context())
		authHeader := gc.GetHeader(("Authorization"))

		if authHeader == "" {
			next.ServeHTTP(w, r)
		} else {
			tokenSplit := strings.Split(authHeader, " ")

			if len(tokenSplit) != 2 {
				next.ServeHTTP(w, r)
			} else {
				token, err := VerifyToken(r.Context(), tokenSplit[1], config)

				if err == nil {
					ctx := context.WithValue(r.Context(), authContextKey, token)

					next.ServeHTTP(w, r.WithContext(ctx))
				} else {
					next.ServeHTTP(w, r)
				}
			}
		}
	})
}
