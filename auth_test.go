package firebasedb

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestSecret(t *testing.T) {
    s := Secret{Token: "password"}
    assert.Equal(t, "password", s.String())
    assert.Error(t, s.Renew())
}
