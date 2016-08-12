package firebasedb

import "errors"

type Authenticator interface {
    String() string
    Renew() error
}

type Secret struct {
    Token string
}

func (s Secret) String() string {
    return s.Token
}

func (s Secret) Renew() error {
    return errors.New("Can't renew a static token")
}