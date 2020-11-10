package golden

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Exchange_request(t *testing.T) {
	// --- When ---
	gld := NewExchange(Open(t, "testdata/request.yaml", nil))

	// --- Then ---
	require.NotNil(t, gld.Request)
	assert.Nil(t, gld.Response)
	assert.Exactly(t, http.MethodPost, gld.Request.Method)
	assert.Exactly(t, "/some/path", gld.Request.Path)
	assert.Exactly(t, "key0=val0&key1=val1", gld.Request.Query)

	exp := []string{
		"Authorization: Bearer token",
		"Content-Type: application/json",
	}
	assert.Exactly(t, exp, gld.Request.Headers)
	assert.Exactly(t, "{\n  \"key2\": \"val2\"\n}\n", gld.Request.Body)
}

func Test_Exchange_response(t *testing.T) {
	// --- When ---
	gld := NewExchange(Open(t, "testdata/response.yaml", nil))

	// --- Then ---
	require.Nil(t, gld.Request)
	assert.NotNil(t, gld.Response)
	assert.Exactly(t, 200, gld.Response.StatusCode)

	exp := []string{
		"Authorization: Bearer token",
		"Content-Type: application/json",
	}
	assert.Exactly(t, exp, gld.Response.Headers)
	assert.Exactly(t, "{ \"key2\": \"val2\" }\n", gld.Response.Body)
}

func Test_Exchange_request_response(t *testing.T) {
	// --- When ---
	gld := NewExchange(Open(t, "testdata/request_response.yaml", nil))

	// --- Then ---
	require.NotNil(t, gld.Request)
	assert.NotNil(t, gld.Response)

	// Request
	assert.Exactly(t, http.MethodPost, gld.Request.Method)
	assert.Exactly(t, "/some/path", gld.Request.Path)
	assert.Exactly(t, "key0=val0&key1=val1", gld.Request.Query)

	exp := []string{
		"Authorization: Bearer token",
		"Content-Type: application/json",
	}
	assert.Exactly(t, exp, gld.Request.Headers)
	assert.Exactly(t, "{\n  \"key2\": \"val2\"\n}\n", gld.Request.Body)

	// Response
	assert.Exactly(t, 200, gld.Response.StatusCode)

	exp = []string{
		"Content-Type: application/json",
	}
	assert.Exactly(t, exp, gld.Response.Headers)
	assert.Exactly(t, "{ \"success\": true }\n", gld.Response.Body)
}

func Test_Exchange_template(t *testing.T) {
	// --- Given ---
	data := map[string]interface{}{
		"val1": 1,
		"val2": "val2",
	}

	// --- When ---
	gld := NewExchange(Open(t, "testdata/request.tpl.yaml", data))

	// --- Then ---
	require.NotNil(t, gld.Request)
	assert.Nil(t, gld.Response)
	assert.Exactly(t, http.MethodPost, gld.Request.Method)
	assert.Exactly(t, "/some/path", gld.Request.Path)
	assert.Exactly(t, "key0=val0&key1=1", gld.Request.Query)

	exp := []string{
		"Authorization: Bearer token",
		"Content-Type: application/json",
	}
	assert.Exactly(t, exp, gld.Request.Headers)
	assert.Exactly(t, "{\n  \"key2\": \"val2\"\n}\n", gld.Request.Body)
}

func Test_Exchange_WriteTo(t *testing.T) {
	// --- Given ---
	gld := NewExchange(Open(t, "testdata/request_response.yaml", nil))
	dst := &bytes.Buffer{}

	// --- When ---
	n, err := gld.WriteTo(dst)

	// --- Then ---
	assert.NoError(t, err)
	assert.Exactly(t, int64(395), n)

	got := NewExchange(t, dst)

	// Request
	assert.Exactly(t, http.MethodPost, got.Request.Method)
	assert.Exactly(t, "/some/path", got.Request.Path)
	assert.Exactly(t, "key0=val0&key1=val1", got.Request.Query)

	exp := []string{
		"Authorization: Bearer token",
		"Content-Type: application/json",
	}
	assert.Exactly(t, exp, got.Request.Headers)
	assert.Exactly(t, "{\n  \"key2\": \"val2\"\n}\n", got.Request.Body)

	// Response
	assert.Exactly(t, 200, got.Response.StatusCode)

	exp = []string{
		"Content-Type: application/json",
	}
	assert.Exactly(t, exp, got.Response.Headers)
	assert.Exactly(t, "{ \"success\": true }\n", got.Response.Body)
}
