package golden

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_helpers_headers2Lines(t *testing.T) {
	// --- Given ---
	hs := http.Header{}
	hs.Add("Authorization", "Bearer token")
	hs.Add("Custom-Header", "val0")
	hs.Add("Custom-Header", "val1")

	// --- When ---
	lns := headers2Lines(t, hs)

	// --- Then ---
	assert.Len(t, lns, 3)
	assert.Exactly(t, "Authorization: Bearer token", lns[0])
	assert.Exactly(t, "Custom-Header: val0", lns[1])
	assert.Exactly(t, "Custom-Header: val1", lns[2])
}

func Test_helpers_headers2Lines_emptyHeaders(t *testing.T) {
	// --- Given ---
	hs := http.Header{}

	// --- When ---
	lns := headers2Lines(t, hs)

	// --- Then ---
	assert.Len(t, lns, 0)
}

func Test_helpers_lines2Headers(t *testing.T) {
	// --- Given ---
	lns := []string{
		"Authorization: Bearer token",
		"Custom-Header: val0",
		"Custom-Header: val1",
	}

	// --- When ---
	hs := lines2Headers(t, lns...)

	// --- Then ---
	assert.Len(t, hs, 2)
	assert.Contains(t, hs, "Authorization")
	assert.Contains(t, hs, "Custom-Header")
	assert.Len(t, hs.Values("Authorization"), 1)
	assert.Len(t, hs.Values("Custom-Header"), 2)
	assert.Exactly(t, "Bearer token", hs.Get("Authorization"))
	assert.Exactly(t, "val0", hs.Values("Custom-Header")[0])
	assert.Exactly(t, "val1", hs.Values("Custom-Header")[1])
}

func Test_helpers_readBody(t *testing.T) {
	// --- Given ---
	body := "{\n    \"key2\": \"val2\"\n}\n"
	req := httptest.NewRequest(
		http.MethodPost,
		"/some/path",
		strings.NewReader(body),
	)

	// --- When ---
	got, rc := readBody(t, req.Body)

	// --- Then ---
	assert.Exactly(t, body, string(got))

	// Make sure you can still read form the request body.
	reqBody, _ := ioutil.ReadAll(rc)
	assert.Exactly(t, body, string(reqBody))
}

func Test_helpers_Map(t *testing.T) {
	// --- When ---
	data := make(Map).Add("key1", "val1").Add("key2", 2)

	// --- Then ---
	exp := map[string]interface{}{
		"key1": "val1",
		"key2": 2,
	}
	assert.Exactly(t, exp, map[string]interface{}(data))
}
