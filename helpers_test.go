package golden

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Helpers_lines2Headers(t *testing.T) {
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

func Test_Helpers_headers2Lines(t *testing.T) {
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

func Test_Helpers_headers2Lines_emptyHeaders(t *testing.T) {
	// --- Given ---
	hs := http.Header{}

	// --- When ---
	lns := headers2Lines(t, hs)

	// --- Then ---
	assert.Len(t, lns, 0)
}

func Test_Helpers_body2Lines(t *testing.T) {
	// --- Given ---
	content := []byte("{\n    \"key2\": \"val2\"\n}\n")
	body := bytes.NewReader(content)
	req := httptest.NewRequest(http.MethodPost, "/some/path", body)

	// --- When ---
	lns := body2Lines(t, req)

	// --- Then ---
	exp := []string{"{", "    \"key2\": \"val2\"", "}", ""}
	assert.Exactly(t, exp, lns)

	// Make sure you can still read form the request body.
	got, _ := ioutil.ReadAll(req.Body)
	assert.Exactly(t, string(content), string(got))
}
