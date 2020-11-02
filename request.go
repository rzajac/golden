package golden

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
)

// Request represents HTTP request backed by a golden file.
type Request struct {
	// Golden file backing this request.
	*Golden

	// Request headers.
	header http.Header
}

// NewRequest returns new instance of Request.
func NewRequest(t T, rdr io.Reader) *Request {
	req := &Request{
		Golden: NewGolden(t, rdr),
	}

	if req.Section(SecReqMethod) == nil {
		req.t.Fatal(errors.New("HTTP request needs request method"))
	}

	if req.Section(SecReqPath) == nil {
		req.t.Fatal(errors.New("HTTP request needs request path"))
	}

	if header := req.Section(SecHeader); header != nil {
		req.header = parseHeaders(req.t, header.lines...)
	} else {
		req.header = make(http.Header)
	}

	return req
}

// AssertRequest asserts request matches the golden file.
// Only the sections that were set are asserted and only the
// headers set in the golden file - in another words request may have
// more headers then the golden file.
func (req *Request) AssertRequest(got *http.Request) {
	req.t.Helper()

	exp := req.Section(SecReqMethod).String()
	if exp != got.Method {
		req.t.Fatalf("expected request method %s got %s", exp, got.Method)
		return
	}

	exp = req.Section(SecReqPath).String()
	if exp != got.URL.Path {
		req.t.Fatalf("expected request path %s got %s", exp, got.URL.Path)
		return
	}

	exp = ""
	if sec := req.Section(SecReqQuery); sec != nil {
		exp = req.Section(SecReqQuery).String()
	}
	if exp != got.URL.RawQuery {
		req.t.Fatalf("expected request query %s got %s", exp, got.URL.RawQuery)
		return
	}

	// Checks only headers set in golden file, got request may have more.
	for key, vv := range req.header {
		g := got.Header.Values(key)
		if !reflect.DeepEqual(vv, g) {
			req.t.Fatalf(
				"expected request header %s values %v got %v",
				key,
				vv,
				g,
			)
			return
		}
	}
}

// Request returns HTTP request matching golden file. It panics on error.
func (req *Request) Request() *http.Request {
	var body io.Reader

	if sec := req.Section(SecBody); sec != nil {
		body = strings.NewReader(sec.String())
	}

	httpReq := httptest.NewRequest(
		req.Section(SecReqMethod).String(),
		req.Section(SecReqPath).String(),
		body,
	)

	if query := req.Section(SecReqQuery); query != nil {
		httpReq.URL.RawQuery = query.String()
	}

	if header := req.Section(SecHeader); header != nil {
		httpReq.Header = parseHeaders(req.t, header.lines...)
	}

	return httpReq
}
