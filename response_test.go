package golden

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rzajac/golden/goldentest"
)

func Test_Response_Assert(t *testing.T) {
	// --- Given ---
	body := `{"key2":"val2"}`
	rsp := &http.Response{
		Header: make(http.Header),
	}
	rsp.StatusCode = 200
	rsp.Header.Add("Authorization", "Bearer token")
	rsp.Header.Add("Content-Type", "application/json")
	rsp.Body = ioutil.NopCloser(strings.NewReader(body))

	// --- When ---
	gld := NewExchange(Open(t, "testdata/response.yaml", nil))

	// --- Then ---
	gld.Response.Assert(rsp)
}

func Test_Response_Assert_HeaderDoesNotMatch(t *testing.T) {
	// --- Given ---
	mck := &goldentest.TMock{}
	mck.On("Helper")
	mck.On(
		"Fatalf",
		"expected response header %s values %v got %v",
		"Authorization",
		[]string{"Bearer token"},
		[]string{"Bearer token 2"},
	)

	body := `{"key2":"val2"}`
	rsp := &http.Response{
		Header: make(http.Header),
	}
	rsp.StatusCode = 200
	rsp.Header.Add("Authorization", "Bearer token 2")
	rsp.Header.Add("Content-Type", "application/json")
	rsp.Body = ioutil.NopCloser(strings.NewReader(body))

	// --- When ---
	gld := NewExchange(Open(mck, "testdata/response.yaml", nil))

	// --- Then ---
	gld.Response.Assert(rsp)
}

func Test_Response_Assert_OnlyDefinedHeadersChecked(t *testing.T) {
	// --- Given ---
	mck := &goldentest.TMock{}
	mck.On("Helper")

	body := `{"key2":"val2"}`
	rsp := &http.Response{
		Header: make(http.Header),
	}
	rsp.StatusCode = 200
	rsp.Header.Add("Authorization", "Bearer token")
	rsp.Header.Add("Content-Type", "application/json")
	rsp.Header.Add("Custom-Header", "custom data")
	rsp.Body = ioutil.NopCloser(strings.NewReader(body))

	// --- When ---
	gld := NewExchange(Open(mck, "testdata/response.yaml", nil))

	// --- Then ---
	gld.Response.Assert(rsp)
}

func Test_Response_Unmarshall(t *testing.T) {
	// --- Given ---
	gld := NewExchange(Open(t, "testdata/response.yaml", nil))

	// --- When ---
	m := make(map[string]string, 1)
	gld.Response.Unmarshall(&m)

	// --- Then ---
	require.Len(t, m, 1)
	require.Contains(t, m, "key2")
	assert.Exactly(t, "val2", m["key2"])
}

func Test_Response_Bytes(t *testing.T) {
	// --- When ---
	gld := NewExchange(Open(t, "testdata/response.yaml", nil))

	// --- Then ---
	exp := []byte("{ \"key2\": \"val2\" }\n")
	assert.Exactly(t, exp, gld.Response.Bytes())
}
