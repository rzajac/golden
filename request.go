package golden

import (
	"encoding/json"
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
		req.header = lines2Headers(req.t, header.lines...)
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
		httpReq.Header = lines2Headers(req.t, header.lines...)
	}

	return httpReq
}

// UnmarshallJSONBody unmarshalls request body to v. Calls Fatal if body
// section does not exist or json.Unmarshal returns error.
func (req *Request) UnmarshallJSONBody(v interface{}) {
	if sec := req.Section(SecBody); sec != nil {
		if err := json.Unmarshal(sec.Bytes(), v); err != nil {
			req.t.Fatal(err)
		}
	}
	req.t.Fatal(errors.New("golden file does not have body"))
}

// RequestSave saves request as golden file.
func RequestSave(t T, w io.Writer, req *http.Request) {
	gld := &Golden{
		sections: make([]*Section, 0, 5),
		t:        nil,
	}

	method := &Section{
		id:    SecReqMethod,
		lines: []string{req.Method},
		mod:   "",
	}
	gld.sections = append(gld.sections, method)

	path := &Section{
		id:    SecReqPath,
		lines: []string{req.URL.Path},
		mod:   "",
	}
	gld.sections = append(gld.sections, path)

	query := &Section{
		id:    SecReqQuery,
		lines: []string{req.URL.RawQuery},
		mod:   "",
	}
	gld.sections = append(gld.sections, query)

	header := &Section{
		id:    SecHeader,
		lines: Headers2Lines(t, req.Header),
		mod:   "",
	}
	gld.sections = append(gld.sections, header)

	body := &Section{
		id:    SecBody,
		lines: body2Lines(t, req),
		mod:   "",
	}
	gld.sections = append(gld.sections, body)

	if _, err := gld.WriteTo(w); err != nil {
		t.Fatal(err)
	}
}
