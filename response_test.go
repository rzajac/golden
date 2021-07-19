package golden

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	. "github.com/rzajac/golden/internal"
)

func Test_Response(t *testing.T) {
	// --- When ---
	gld := NewResponse(Open(t, "testdata/response.yaml", nil))

	// --- Then ---
	assert.Exactly(t, 200, gld.StatusCode)

	exp := []string{
		"Authorization: Bearer token",
		"Content-Type: application/json",
	}
	assert.Exactly(t, exp, gld.Headers)
	assert.Exactly(t, "{ \"key2\": \"val2\" }\n", gld.Body)

	assert.Exactly(t, "val1", gld.Meta["key1"])
	assert.Exactly(t, 123, gld.Meta["key2"])
	assert.Exactly(t, 12.3, gld.Meta["key3"])

	expDate := time.Date(2021, 2, 28, 10, 24, 25, 123000000, time.UTC)
	assert.Exactly(t, expDate, gld.Meta["key4"].(time.Time))
}

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
	gld := NewResponse(Open(t, "testdata/response.yaml", nil))

	// --- Then ---
	gld.Assert(rsp)
}

func Test_Response_Assert_diff(t *testing.T) {
	// --- Given ---
	body := `{"key1":"val1","key3":"val3","key4":"val4"}`
	rsp := &http.Response{
		Header: make(http.Header),
	}
	rsp.StatusCode = 200
	rsp.Body = ioutil.NopCloser(strings.NewReader(body))

	mck := &TMock{}
	mck.On("Helper")
	mck.On("Fatal", mock.AnythingOfType("string"))
	gld := NewResponse(Open(mck, "testdata/response2.yaml", nil))

	// --- When ---
	gld.Assert(rsp)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Response_Assert_HeaderDoesNotMatch(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
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
	gld := NewResponse(Open(mck, "testdata/response.yaml", nil))

	// --- Then ---
	gld.Assert(rsp)
}

func Test_Response_Assert_OnlyDefinedHeadersChecked(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
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
	gld := NewResponse(Open(mck, "testdata/response.yaml", nil))

	// --- Then ---
	gld.Assert(rsp)
}

func Test_Response_Unmarshal(t *testing.T) {
	// --- Given ---
	gld := NewResponse(Open(t, "testdata/response.yaml", nil))

	// --- When ---
	m := make(map[string]string, 1)
	gld.Unmarshal(&m)

	// --- Then ---
	require.Len(t, m, 1)
	require.Contains(t, m, "key2")
	assert.Exactly(t, "val2", m["key2"])
}

func Test_Response_Bytes(t *testing.T) {
	// --- When ---
	gld := NewResponse(Open(t, "testdata/response.yaml", nil))

	// --- Then ---
	exp := []byte("{ \"key2\": \"val2\" }\n")
	assert.Exactly(t, exp, gld.Bytes())
}
