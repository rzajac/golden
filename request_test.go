package golden

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/rzajac/golden/internal"
)

func Test_Request(t *testing.T) {
	// --- When ---
	gld := NewRequest(Open(t, "testdata/request.yaml", nil))

	// --- Then ---
	require.NotNil(t, gld)
	assert.Exactly(t, http.MethodPost, gld.Method)
	assert.Exactly(t, "/some/path", gld.Path)
	assert.Exactly(t, "key0=val0&key1=val1", gld.Query)

	exp := []string{
		"Authorization: Bearer token",
		"Content-Type: application/json",
	}
	assert.Exactly(t, exp, gld.Headers)
	assert.Exactly(t, "{\n  \"key2\": \"val2\"\n}\n", gld.Body)
	assert.Exactly(t, "https", gld.Scheme)

	assert.Exactly(t, "val1", gld.Meta["key1"])
	assert.Exactly(t, 123, gld.Meta["key2"])
	assert.Exactly(t, 12.3, gld.Meta["key3"])

	expDate := time.Date(2021, 2, 28, 10, 24, 25, 123000000, time.UTC)
	assert.Exactly(t, expDate, gld.Meta["key4"].(time.Time))
}

func Test_Request_Assert(t *testing.T) {
	// --- Given ---
	req := httptest.NewRequest(
		http.MethodPost,
		"/some/path",
		strings.NewReader(`{"key2":"val2"}`),
	)
	req.Header.Add("Authorization", "Bearer token")
	req.Header.Add("Content-Type", "application/json")
	req.URL.RawQuery = "key0=val0&key1=val1"

	// --- When ---
	gld := NewRequest(Open(t, "testdata/request.yaml", nil))

	// --- Then ---
	gld.Assert(req)
}

func Test_Request_Assert_MethodDoesNotMatch(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On("Fatalf", "expected request method %s got %s", "POST", "GET")

	req := httptest.NewRequest(http.MethodGet, "/some/path", nil)

	// --- When ---
	gld := NewRequest(Open(mck, "testdata/request.yaml", nil))

	// --- Then ---
	gld.Assert(req)
}

func Test_Request_Assert_PathDoesNotMatch(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On(
		"Fatalf",
		"expected request path %s got %s",
		"/some/path",
		"/other/path",
	)

	req := httptest.NewRequest(http.MethodPost, "/other/path", nil)

	// --- When ---
	gld := NewRequest(Open(mck, "testdata/request.yaml", nil))

	// --- Then ---
	gld.Assert(req)
}

func Test_Request_Assert_QueryDoesNotMatch(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On(
		"Fatalf",
		"expected request query %s got %s",
		"key0=val0&key1=val1",
		"key0=val0",
	)

	req := httptest.NewRequest(http.MethodPost, "/some/path", nil)
	req.URL.RawQuery = "key0=val0"

	// --- When ---
	gld := NewRequest(Open(mck, "testdata/request.yaml", nil))

	// --- Then ---
	gld.Assert(req)
}

func Test_Request_Assert_HeaderDoesNotMatch(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On(
		"Fatalf",
		"expected request header %s values %v got %v",
		"Authorization",
		[]string{"Bearer token"},
		[]string{"Bearer token2"},
	)

	req := httptest.NewRequest(http.MethodPost, "/some/path", nil)
	req.URL.RawQuery = "key0=val0&key1=val1"
	req.Header.Add("Authorization", "Bearer token2")

	// --- When ---
	gld := NewRequest(Open(mck, "testdata/request.yaml", nil))

	// --- Then ---
	gld.Assert(req)
}

func Test_Request_Assert_OnlyDefinedHeadersChecked(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")

	body := strings.NewReader(`{"key2":"val2"}`)
	req := httptest.NewRequest(http.MethodPost, "/some/path", body)
	req.URL.RawQuery = "key0=val0&key1=val1"
	req.Header.Add("Authorization", "Bearer token")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Custom-Header", "custom data")

	// --- When ---
	gld := NewRequest(Open(mck, "testdata/request.yaml", nil))

	// --- Then ---
	gld.Assert(req)
}

func Test_Request_Request(t *testing.T) {
	// --- Given ---
	gld := NewRequest(Open(t, "testdata/request.yaml", nil))

	// --- When ---
	got := gld.Request()

	// --- Then ---
	assert.Exactly(t, http.MethodPost, got.Method)
	assert.Exactly(t, "/some/path", got.URL.Path)
	assert.Exactly(t, "key0=val0&key1=val1", got.URL.RawQuery)
	require.Len(t, got.Header, 2)
	require.Contains(t, got.Header, "Authorization")
	require.Contains(t, got.Header, "Content-Type")
	require.Len(t, got.Header.Values("Authorization"), 1)
	require.Len(t, got.Header.Values("Content-Type"), 1)
	assert.Exactly(t, "Bearer token", got.Header.Values("Authorization")[0])
	assert.Exactly(t, "application/json", got.Header.Values("Content-Type")[0])
}

func Test_Request_Unmarshal(t *testing.T) {
	// --- Given ---
	gld := NewRequest(Open(t, "testdata/request.yaml", nil))

	// --- When ---
	m := make(map[string]string, 1)
	gld.Unmarshal(&m)

	// --- Then ---
	require.Len(t, m, 1)
	require.Contains(t, m, "key2")
	assert.Exactly(t, "val2", m["key2"])
}

func Test_Request_BindQuery(t *testing.T) {
	// --- Given ---
	gld := NewRequest(Open(t, "testdata/request.yaml", nil))

	type T1 struct {
		Key0 string `form:"key0"`
		Key1 string `form:"key1"`
	}

	// --- When ---
	t1 := &T1{}
	gld.BindQuery("form", t1)

	// --- Then ---
	assert.Exactly(t, "val0", t1.Key0)
	assert.Exactly(t, "val1", t1.Key1)
}

func Test_Request_Bytes(t *testing.T) {
	// --- When ---
	gld := NewRequest(Open(t, "testdata/request.yaml", nil))

	// --- Then ---
	exp := []byte("{\n  \"key2\": \"val2\"\n}\n")
	assert.Exactly(t, exp, gld.Bytes())
}
