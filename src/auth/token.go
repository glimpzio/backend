package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Token struct {
	AuthId string
	Expiry int64
}

var ACCESS_TOKEN_COOKIE = "ACCESS_TOKEN"
var REFRESH_TOKEN_COOKIE = "REFRESH_TOKEN"

type token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// Exchange an auth code
func ExchangeAuthCode(auth0Config *Auth0Config, code string) (*token, error) {
	payload := strings.NewReader(fmt.Sprintf("grant_type=authorization_code&client_id=%s&client_secret=%s&code=%s&redirect_uri=%s", auth0Config.Auth0ClientId, auth0Config.Auth0ClientSecret, code, auth0Config.Auth0RedirectUrl))
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/oauth/token", auth0Config.Auth0Domain), payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	tkn := &token{}
	if err := json.Unmarshal(body, tkn); err != nil {
		return nil, err
	}

	return tkn, nil
}

// Refresh token
func RefreshToken(auth0Config *Auth0Config, refreshToken string) (*token, error) {
	payload := strings.NewReader(fmt.Sprintf("grant_type=refresh_token&client_id=%s&client_secret=%s&refresh_token=%s", auth0Config.Auth0ClientId, auth0Config.Auth0ClientSecret, refreshToken))
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/oauth/token", auth0Config.Auth0Domain), payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	tkn := &token{}
	if err := json.Unmarshal(body, tkn); err != nil {
		return nil, err
	}

	return tkn, nil
}
