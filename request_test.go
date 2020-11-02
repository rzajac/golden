package golden

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rzajac/golden/goldentest"
)

func Test_Request_Request(t *testing.T) {
	// --- Given ---
	req := NewRequest(t, Open(t, "testdata/request_basic.gold"))

	// --- When ---
	httpReq := req.Request()

	// --- Then ---
	assert.Exactly(t, http.MethodPost, httpReq.Method)
	assert.Exactly(t, "/some/path", httpReq.URL.Path)
	assert.Exactly(t, "key0=val0&key1=val1", httpReq.URL.RawQuery)

	exp := map[string][]string{
		"Authorization": {"Bearer token"},
	}
	assert.Exactly(t, http.Header(exp), httpReq.Header)

	b, _ := ioutil.ReadAll(httpReq.Body)
	assert.Exactly(t, `{"key2": "val2"}`, string(b))
}

func Test_Request_AssertRequest(t *testing.T) {
	// --- Given ---
	body := `{"key2": "val2"}`
	req := httptest.NewRequest(
		http.MethodPost,
		"/some/path",
		strings.NewReader(body),
	)
	req.Header.Add("Authorization", "Bearer token")
	req.Header.Add("Host", "localhost")
	req.URL.RawQuery = "key0=val0&key1=val1"

	// --- Then ---
	NewRequest(t, Open(t, "testdata/request_basic.gold")).AssertRequest(req)
}

func Test_Request_AssertRequest_MethodDoesNotMatch(t *testing.T) {
	// --- Given ---
	mck := &goldentest.TMock{}
	mck.On("Helper")
	mck.On("Fatalf", "expected request method %s got %s", "POST", "GET")

	req := httptest.NewRequest(http.MethodGet, "/some/path", nil)

	// --- When ---
	NewRequest(mck, Open(mck, "testdata/request_basic.gold")).AssertRequest(req)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Request_AssertRequest_PathDoesNotMatch(t *testing.T) {
	// --- Given ---
	mck := &goldentest.TMock{}
	mck.On("Helper")
	mck.On(
		"Fatalf",
		"expected request path %s got %s",
		"/some/path",
		"/other/path",
	)

	req := httptest.NewRequest(http.MethodPost, "/other/path", nil)

	// --- When ---
	NewRequest(mck, Open(mck, "testdata/request_basic.gold")).AssertRequest(req)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Request_AssertRequest_QueryDoesNotMatch(t *testing.T) {
	// --- Given ---
	mck := &goldentest.TMock{}
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
	NewRequest(mck, Open(mck, "testdata/request_basic.gold")).AssertRequest(req)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Request_AssertRequest_HeaderDoesNotMatch(t *testing.T) {
	// --- Given ---
	mck := &goldentest.TMock{}
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
	NewRequest(mck, Open(mck, "testdata/request_basic.gold")).AssertRequest(req)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Request_AssertRequest_OnlyDefinedHeadersChecked(t *testing.T) {
	// --- Given ---
	mck := &goldentest.TMock{}
	mck.On("Helper")

	req := httptest.NewRequest(http.MethodPost, "/some/path", nil)
	req.URL.RawQuery = "key0=val0&key1=val1"
	req.Header.Add("Authorization", "Bearer token")
	req.Header.Add("Custom-Header", "custom data")

	// --- When ---
	NewRequest(mck, Open(mck, "testdata/request_basic.gold")).AssertRequest(req)

	// --- Then ---
	mck.AssertNotCalled(t, "Fatalf")
}
