package infrastructure

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSYRAuth(t *testing.T) {
	urlauth := "http://syr.ru/auth/"
	tokenfield := "message"
	tokenheader := "yadro"
	login := "test"
	password := "test"

	t.Run("valid New", func(t *testing.T) {
		a, err := NewSYRAuth(
			urlauth,
			tokenfield,
			tokenheader,
			login,
			password)
		assert.NotNil(t, a)
		assert.Nil(t, err)
	})
	t.Run("invalid url", func(t *testing.T) {
		a, err := NewSYRAuth(
			"",
			tokenfield,
			tokenheader,
			login,
			password)
		assert.Nil(t, a)
		assert.NotNil(t, err)
	})
	t.Run("invalid tokenfield", func(t *testing.T) {
		a, err := NewSYRAuth(
			urlauth,
			"",
			tokenheader,
			login,
			password)
		assert.Nil(t, a)
		assert.NotNil(t, err)
	})
	t.Run("invalid tokenheader", func(t *testing.T) {
		a, err := NewSYRAuth(
			urlauth,
			tokenfield,
			"",
			login,
			password)
		assert.Nil(t, a)
		assert.NotNil(t, err)
	})
	t.Run("invalid login", func(t *testing.T) {
		a, err := NewSYRAuth(
			urlauth,
			tokenfield,
			tokenheader,
			"",
			password)
		assert.Nil(t, a)
		assert.NotNil(t, err)
	})
	t.Run("invalid password", func(t *testing.T) {
		a, err := NewSYRAuth(
			urlauth,
			tokenfield,
			tokenheader,
			login,
			"")
		assert.Nil(t, a)
		assert.NotNil(t, err)
	})
}
func TestAuthorize(t *testing.T) {
	tokenfield := "message"
	token := "0123456789"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{\"%s\": \"%s\"}\n", tokenfield, token)
	}))
	defer ts.Close()
	urlauth := ts.URL
	tokenheader := "yadro"
	login := "test"
	password := "test"
	a, _ := NewSYRAuth(
		urlauth,
		tokenfield,
		tokenheader,
		login,
		password)
	assert.NotNil(t, a)
	ctx := context.Background()
	err := a.authorize(ctx)
	assert.Nil(t, err)
	assert.Equal(t, tokenheader+token, a.token())
}
