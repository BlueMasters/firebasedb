package firebasedb

import "errors"

// Authenticator is the interface used to add authentication data to the requests. The String() method
// returns the current token and Renew() is called if the current token has expired.
type Authenticator interface {
	String() string
	ParamName() string // usually "auth" ou "access_token"
	Renew() error
}

// Secret implements the Authenticator interface and is used with static Database secret.
type Secret struct {
	Token string
}

// String returns the static Database secret
func (s Secret) String() string {
	return s.Token
}

func (s Secret) ParamName() string {
	return "auth"
}

// Renew is not allowed for static secret and thus always returns an error.
func (s Secret) Renew() error {
	return errors.New("Can't renew a static token")
}
