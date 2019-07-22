package provider

type Providers struct {
	Google Google `group:"Google Provider" namespace:"google" env-namespace:"GOOGLE"`
	OAuth2 OAuth2 `group:"OAuth2 Provider" namespace:"oauth2" env-namespace:"OAUTH2"`
}

type Token struct {
	Token string `json:"access_token"`
}

type User struct {
	Email string `json:"email"`
}
