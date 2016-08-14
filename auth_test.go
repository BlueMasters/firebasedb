package firebasedb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSecret(t *testing.T) {
	s := Secret{Token: "password"}
	assert.Equal(t, "password", s.String())
	assert.Error(t, s.Renew())
}
