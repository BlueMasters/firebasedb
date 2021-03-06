package firebasedb

import (
	"github.com/cenkalti/backoff"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPassKeepAlive(t *testing.T) {
	r1 := NewReference("https://domain.com/")
	assert.NoError(t, r1.Error)
	assert.False(t, r1.passKeepAlive)

	r2 := r1.PassKeepAlive(true)
	assert.False(t, r1.passKeepAlive)
	assert.True(t, r2.passKeepAlive)

	r3 := r2.PassKeepAlive(false)
	assert.False(t, r1.passKeepAlive)
	assert.True(t, r2.passKeepAlive)
	assert.False(t, r3.passKeepAlive)
}

func TestPassRetry(t *testing.T) {
	r1 := NewReference("https://domain.com/")
	assert.NoError(t, r1.Error)
	assert.Nil(t, r1.retry)

	r2 := r1.Retry(backoff.NewExponentialBackOff())
	assert.Nil(t, r1.retry)
	assert.IsType(t, backoff.NewExponentialBackOff(), r2.retry)

	r3 := r2.Retry(nil)
	assert.Nil(t, r1.retry)
	assert.NotNil(t, r2.retry)
	assert.Nil(t, r3.retry)
}

func TestRules(t *testing.T) {
	r1 := NewReference("https://domain.com/")
	assert.NoError(t, r1.Error)
	assert.Equal(t, "/", r1.url.Path)

	r2 := r1.Rules()
	assert.Equal(t, "/", r1.url.Path)
	assert.Equal(t, "/.settings/rules", r2.url.Path)
}
