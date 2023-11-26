package auth

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
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
	Auth0AudienceApi  string
	Auth0RedirectUrl  string
}

// Verify a token
func VerifyToken(ctx context.Context, accessToken string, config *Auth0Config) (*Token, error) {
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

	return &Token{AuthId: claims.RegisteredClaims.Subject, Expiry: claims.RegisteredClaims.Expiry}, nil
}

// Apply middleware
func ApplyMiddleware(logger *misc.Logger, next http.Handler, config *Auth0Config, domain string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gc, err := misc.GinContextFromContext(r.Context())

		if err != nil {
			next.ServeHTTP(w, r)
		} else {
			accessToken := ""
			accessTokenCookie, err := gc.Cookie(ACCESS_TOKEN_COOKIE)
			if err == nil {
				accessToken = accessTokenCookie
			}

			refreshToken := ""
			refreshTokenCookie, err := gc.Cookie(REFRESH_TOKEN_COOKIE)
			if err == nil {
				refreshToken = refreshTokenCookie
			}

			if accessToken == "" && refreshToken == "" {
				next.ServeHTTP(w, r)
			} else {
				token, err := VerifyToken(r.Context(), accessToken, config)

				if err != nil {
					tkn, err := RefreshToken(config, refreshToken)

					if err != nil {
						next.ServeHTTP(w, r)
					} else {
						token, err := VerifyToken(r.Context(), tkn.AccessToken, config)

						if err != nil {
							next.ServeHTTP(w, r)
						} else {
							gc.SetCookie(ACCESS_TOKEN_COOKIE, tkn.AccessToken, tkn.ExpiresIn, "/", domain, true, true)
							gc.SetCookie(REFRESH_TOKEN_COOKIE, tkn.RefreshToken, math.MaxInt, "/", domain, true, true)

							ctx := context.WithValue(r.Context(), authContextKey, token)

							next.ServeHTTP(w, r.WithContext(ctx))
						}
					}
				} else {
					ctx := context.WithValue(r.Context(), authContextKey, token)

					next.ServeHTTP(w, r.WithContext(ctx))
				}
			}
		}

	})
}
