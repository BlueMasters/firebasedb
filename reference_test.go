package firebasedb

import (
    "testing"
    "github.com/stretchr/testify/assert"
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
    assert.False(t, r1.retry)

    r2 := r1.Retry(true)
    assert.False(t, r1.retry)
    assert.True(t, r2.retry)

    r3 := r2.Retry(false)
    assert.False(t, r1.retry)
    assert.True(t, r2.retry)
    assert.False(t, r3.retry)
}

func TestRules(t *testing.T) {
    r1 := NewReference("https://domain.com/")
    assert.NoError(t, r1.Error)
    assert.Equal(t, "/", r1.url.Path)

    r2 := r1.Rules()
    assert.Equal(t, "/", r1.url.Path)
    assert.Equal(t, "/.settings/rules", r2.url.Path)
}