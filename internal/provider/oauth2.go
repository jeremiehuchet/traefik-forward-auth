package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type OAuth2 struct {
	ClientId     string `long:"client-id" env:"CLIENT_ID" description:"Client ID"`
	ClientSecret string `long:"client-secret" env:"CLIENT_SECRET" description:"Client Secret" json:"-"`
	Scope        string `long:"scope" env:"SCOPE" description:"OAuth2 scope"`
	Prompt       string `long:"prompt" env:"PROMPT" description:"Space separated list of OpenID prompt options"`

	LoginURL string `long:"login-url" env:"LOGIN_URL" description:"OAuth2 login endpoint"`
	TokenURL string `long:"token-url" env:"TOKEN_URL" description:"OAuth2 token endpoint"`
	UserURL  string `long:"user-url" env:"USER_URL" description:"User info endpoint"`
}

func (g *OAuth2) GetLoginURL(redirectUri, state string) string {
	q := url.Values{}
	q.Set("client_id", g.ClientId)
	q.Set("response_type", "code")
	q.Set("scope", g.Scope)
	if g.Prompt != "" {
		q.Set("prompt", g.Prompt)
	}
	q.Set("redirect_uri", redirectUri)
	q.Set("state", state)

	u, _ := url.Parse(g.LoginURL)
	u.RawQuery = q.Encode()

	return u.String()
}

func (g *OAuth2) ExchangeCode(redirectUri, code string) (string, error) {
	form := url.Values{}
	form.Set("client_id", g.ClientId)
	form.Set("client_secret", g.ClientSecret)
	form.Set("grant_type", "authorization_code")
	form.Set("redirect_uri", redirectUri)
	form.Set("code", code)

	res, err := http.PostForm(g.TokenURL, form)
	if err != nil {
		return "", err
	}

	var token Token
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&token)

	return token.Token, err
}

func (g *OAuth2) GetUser(token string) (User, error) {
	var user User

	client := &http.Client{}
	req, err := http.NewRequest("GET", g.UserURL, nil)
	if err != nil {
		return user, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err := client.Do(req)
	if err != nil {
		return user, err
	}

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&user)

	return user, err
}
