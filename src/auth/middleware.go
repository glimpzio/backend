package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/glimpzio/backend/misc"
)

type contextKey string

const authContextKey contextKey = "authKey"

type AuthToken struct {
	AuthId string
}

type Middleware struct {
	Token *AuthToken
}

// Middleware resolver
func GetMiddleware(ctx context.Context) *Middleware {
	out := &Middleware{Token: nil}

	token := ctx.Value(authContextKey)

	if token != nil {
		out.Token = token.(*AuthToken)
	}

	return out
}

type Auth0Config struct {
	Auth0Domain       string
	Auth0ClientId     string
	Auth0ClientSecret string
	Auth0AudienceApi  string
	Auth0RedirectUrl  string
}

// Verify a token
func VerifyToken(ctx context.Context, accessToken string, config *Auth0Config) (*AuthToken, error) {
	issuerUrl, err := url.Parse(fmt.Sprintf("https://%s/", config.Auth0Domain))
	if err != nil {
		return nil, err
	}

	provider := jwks.NewCachingProvider(issuerUrl, time.Duration(time.Duration.Minutes(5)))

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerUrl.String(),
		[]string{config.Auth0AudienceApi},
		validator.WithAllowedClockSkew(time.Duration(time.Duration.Minutes(1))),
	)
	if err != nil {
		return nil, err
	}

	validated, err := jwtValidator.ValidateToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	claims := validated.(*validator.ValidatedClaims)

	return &AuthToken{AuthId: claims.RegisteredClaims.Subject}, nil
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
