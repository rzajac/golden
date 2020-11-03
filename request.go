package golden

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
)

// Request represents HTTP request backed by a golden file.
type Request struct {
	Method  string   `yaml:"method"`
	Path    string   `yaml:"path"`
	Query   string   `yaml:"query"`
	Headers []string `yaml:"headers"`
	Body    string   `yaml:"body"`

	headers http.Header // Request headers.
	t       T           // Test manager.
}

// Validate validates request loaded from golden file.
func (req *Request) Validate() {
	if req.Method == "" {
		req.t.Fatal(errors.New("HTTP request needs request method"))
	}

	if req.Path == "" {
		req.t.Fatal(errors.New("HTTP request needs request path"))
	}

	if len(req.Headers) > 0 {
		req.headers = lines2Headers(req.t, req.Headers...)
	} else {
		req.headers = make(http.Header)
	}
}

// Assert asserts request matches the golden file.
func (req *Request) Assert(got *http.Request) {
	req.t.Helper()

	if req.Method != got.Method {
		req.t.Fatalf("expected request method %s got %s", req.Method, got.Method)
		return
	}

	if req.Path != got.URL.Path {
		req.t.Fatalf("expected request path %s got %s", req.Path, got.URL.Path)
		return
	}

	if req.Query != got.URL.RawQuery {
		req.t.Fatalf("expected request query %s got %s", req.Query, got.URL.RawQuery)
		return
	}

	// Checks only headers set in golden file, got request may have more.
	for key, vv := range req.headers {
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

	body, rc := readBody(req.t, got.Body)
	defer func() { got.Body = rc }()
	if req.Body != body {
		req.t.Fatalf(
			"expected request body to match want\n %s\ngot\n%s",
			req.Body,
			body,
		)
		return
	}
}

// Request returns HTTP request matching golden file. It panics on error.
func (req *Request) Request() *http.Request {
	httpReq := httptest.NewRequest(
		req.Method,
		req.Path,
		strings.NewReader(req.Body),
	)
	httpReq.URL.RawQuery = req.Query
	httpReq.Header = lines2Headers(req.t, req.Headers...)

	return httpReq
}

// UnmarshallJSONBody unmarshalls request body to v. Calls Fatal if body
// section does not exist or json.Unmarshal returns error.
func (req *Request) UnmarshallJSONBody(v interface{}) {
	if req.Body != "" {
		if err := json.Unmarshal([]byte(req.Body), v); err != nil {
			req.t.Fatal(err)
		}
		return
	}
	req.t.Fatal(errors.New("golden file does not have body"))
}
