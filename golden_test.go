package golden

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RequestResponse_request(t *testing.T) {
	// --- When ---
	gld := RequestResponse(t, Open(t, "testdata/request.yaml"))

	// --- Then ---
	require.NotNil(t, gld.Request)
	assert.Nil(t, gld.Response)
	assert.Exactly(t, "POST", gld.Request.Method)
	assert.Exactly(t, "/some/path", gld.Request.Path)
	assert.Exactly(t, "key0=val0&key1=val1", gld.Request.Query)

	exp := []string{
		"Authorization: Bearer token",
		"Content-Type: application/json",
	}
	assert.Exactly(t, exp, gld.Request.Headers)
	assert.Exactly(t, "{\n  \"key2\": \"val2\"\n}\n", gld.Request.Body)
}

func Test_RequestResponse_response(t *testing.T) {
	// --- When ---
	gld := RequestResponse(t, Open(t, "testdata/response.yaml"))

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

func Test_RequestResponse_request_response(t *testing.T) {
	// --- When ---
	gld := RequestResponse(t, Open(t, "testdata/request_response.yaml"))

	// --- Then ---
	require.NotNil(t, gld.Request)
	assert.NotNil(t, gld.Response)

	// Request
	assert.Exactly(t, "POST", gld.Request.Method)
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

func Test_RequestResponse_template(t *testing.T) {
	// --- Given ---
	data := map[string]interface{}{
		"val1": 1,
		"val2": "val2",
	}

	// --- When ---
	gld := RequestResponse(t, OpenTpl(t, "testdata/request.tpl.yaml", data))

	// --- Then ---
	require.NotNil(t, gld.Request)
	assert.Nil(t, gld.Response)
	assert.Exactly(t, "POST", gld.Request.Method)
	assert.Exactly(t, "/some/path", gld.Request.Path)
	assert.Exactly(t, "key0=val0&key1=1", gld.Request.Query)

	exp := []string{
		"Authorization: Bearer token",
		"Content-Type: application/json",
	}
	assert.Exactly(t, exp, gld.Request.Headers)
	assert.Exactly(t, "{\n  \"key2\": \"val2\"\n}\n", gld.Request.Body)
}

func Test_RequestResponse_WriteTo(t *testing.T) {
	// --- Given ---
	gld := RequestResponse(t, Open(t, "testdata/request_response.yaml"))
	dst := &bytes.Buffer{}

	// --- When ---
	n, err := gld.WriteTo(dst)

	// --- Then ---
	assert.NoError(t, err)
	assert.Exactly(t, int64(395), n)

	got := RequestResponse(t, dst.Bytes())

	// Request
	assert.Exactly(t, "POST", got.Request.Method)
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
