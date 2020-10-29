package golden

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Golden_basic(t *testing.T) {
	// --- When ---
	gld := New(t, "testdata/basic.txt")

	// --- Then ---
	assert.Exactly(t, http.MethodPost, gld.Verb)
	assert.True(t, gld.VerbSet)
	assert.Exactly(t, "/some/path", gld.Path)
	assert.True(t, gld.PathSet)
	assert.Exactly(t, "key0=val0&key1=val1", gld.Query.Encode())
	assert.True(t, gld.QuerySet)

	assert.Len(t, gld.Headers, 1)
	assert.True(t, gld.HeadersSet)
	assert.Contains(t, gld.Headers, "Authorization")
	assert.Exactly(t, []string{"Bearer token"}, gld.Headers.Values("Authorization"))

	assert.Exactly(t, `{"key2": "val2"}`, gld.Body.String())
}

func Test_Golden_multi_header(t *testing.T) {
	// --- When ---
	gld := New(t, "testdata/multi_header.txt")

	// --- Then ---
	assert.Len(t, gld.Headers, 2)
	assert.Contains(t, gld.Headers, "Authorization")
	assert.Exactly(t, []string{"Bearer token"}, gld.Headers.Values("Authorization"))
	assert.Exactly(t, []string{"application/json"}, gld.Headers.Values("Content-Type"))
}

func Test_Golden_multi_line_body(t *testing.T) {
	// --- When ---
	gld := New(t, "testdata/body_multi_line_json.txt")

	// --- Then ---
	assert.JSONEq(t, `{"key2": "val2"}`, gld.Body.String())
}

func Test_Golden_multi_line_text(t *testing.T) {
	// --- When ---
	gld := New(t, "testdata/body_multi_line_text.txt")

	// --- Then ---
	assert.Exactly(t, "line 0\nline 1\n", gld.Body.String())
}

func Test_Golden_file_open_error(t *testing.T) {
	// --- Given ---
	var called bool
	opt0 := func(gld *Golden) {
		gld.fatal = func(args ...interface{}) {
			called = true
		}
	}

	// --- When ---
	New(t, "invalid/path", opt0)

	// --- Then ---
	assert.True(t, called)
}

func Test_Golden_query_parse_error(t *testing.T) {
	// --- Given ---
	var called bool
	opt0 := func(gld *Golden) {
		gld.fatal = func(args ...interface{}) {
			called = true
		}
	}

	// --- When ---
	New(t, "testdata/invalid_query.txt", opt0)

	// --- Then ---
	assert.True(t, called)
}

func Test_Golden_new_http_request(t *testing.T) {
	// --- Given ---
	gld := New(t, "testdata/basic.txt")

	// --- When ---
	req := gld.Request()

	// --- Then ---
	assert.Exactly(t, http.MethodPost, req.Method)
	assert.Exactly(t, "/some/path", req.URL.Path)
	assert.Exactly(t, "key0=val0&key1=val1", req.URL.RawQuery)

	exp := map[string][]string{
		"Authorization": {"Bearer token"},
	}
	assert.Exactly(t, http.Header(exp), req.Header)

	b, _ := ioutil.ReadAll(req.Body)
	assert.Exactly(t, `{"key2": "val2"}`, string(b))
}

func Test_Golden_assert_http_request(t *testing.T) {
	// --- When ---
	gld := New(t, "testdata/basic.txt")

	// --- Then ---
	gld.AssertRequest(gld.Request())
}

func Test_SaveRequest(t *testing.T) {
	// --- Given ---
	dst := t.TempDir()

	gld := New(t, "testdata/basic.txt")
	req := gld.Request()

	// --- When ---
	err := SaveRequest(dst+"/basic.txt", req, gld.Comments...)

	// --- Then ---
	require.NoError(t, err)

	exp, err := ioutil.ReadFile("testdata/basic.txt")
	require.NoError(t, err)
	got, err := ioutil.ReadFile(dst + "/basic.txt")
	require.NoError(t, err)
	assert.Exactly(t, string(exp), string(append(got, []byte("\n")...)))
}
