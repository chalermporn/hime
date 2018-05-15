package hime_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/acoshift/hime"
	"github.com/stretchr/testify/assert"
)

func TestGlobal(t *testing.T) {
	app := hime.New()

	assert.Nil(t, app.Global("key1"))
	assert.Nil(t, app.Global("key2"))

	app.
		Globals(hime.Globals{
			"key1": "value1",
			"key2": "value2",
		}).
		Handler(hime.H(func(ctx hime.Context) hime.Result {
			assert.Equal(t, "value1", ctx.Global("key1"))
			assert.Equal(t, "value2", ctx.Global("key2"))
			assert.Nil(t, ctx.Global("invalid"))
			return ctx.Nothing()
		}))

	assert.Equal(t, "value1", app.Global("key1"))
	assert.Equal(t, "value2", app.Global("key2"))
	assert.Nil(t, app.Global("invalid"))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
}