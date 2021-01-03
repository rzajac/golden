package golden

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// Request represents golden file for HTTP request.
type Request struct {
	Method   string   `yaml:"method"`
	Path     string   `yaml:"path"`
	Query    string   `yaml:"query"`
	Headers  []string `yaml:"headers"`
	BodyType string   `yaml:"bodyType"`
	Body     string   `yaml:"body"`

	// Request headers parsed from Headers field during validation.
	headers http.Header

	// Test manager.
	t T
}

// NewRequest returns new instance of Request.
func NewRequest(t T, r io.Reader) *Request {
	t.Helper()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
		return nil
	}

	req := &Request{}
	if err := yaml.Unmarshal(data, req); err != nil {
		t.Fatal(err)
		return nil
	}
	req.t = t
	req.validate()

	return req
}

// validate validates request loaded from golden file.
func (req *Request) validate() {
	if req.Method == "" {
		req.t.Fatal(errors.New("HTTP request needs request method"))
		return
	}

	if req.Path == "" {
		req.t.Fatal(errors.New("HTTP request needs request path"))
		return
	}

	if len(req.Headers) > 0 {
		req.headers = lines2Headers(req.t, req.Headers...)
	} else {
		req.headers = make(http.Header)
	}
}

// Assert asserts request matches the golden file.
//
// All headers defined in the golden file must match exactly but passed
// request may have more headers then defined in the golden file.
//
// To compare request bodies the method best for defined body type is used.
// For example when comparing JSON bodies both byte slices don't have to be
// identical but they must represent the same data.
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

	var equal bool
	switch req.BodyType {
	case TypeJSON:
		equal = assertJSONEqual(req.t, req.Bytes(), body)
	case TypeText:
		equal = bytes.Equal(req.Bytes(), body)
	default:
		equal = bytes.Equal(req.Bytes(), body)
	}

	if !equal {
		req.t.Fatalf(
			"expected request body to match want\n %s\ngot\n%s",
			req.Body,
			body,
		)
		return
	}
}

// Request returns HTTP request represented by the golden file. It panics
// on error.
func (req *Request) Request() *http.Request {
	req.t.Helper()
	httpReq := httptest.NewRequest(
		req.Method,
		req.Path,
		strings.NewReader(req.Body),
	)
	httpReq.URL.RawQuery = req.Query
	httpReq.Header = lines2Headers(req.t, req.Headers...)
	return httpReq
}

// Unmarshal unmarshalls request body to v based on body type. When
// body type is set to text v can be pointer to sting or byte slice (with
// enough space to fit body). Calls Fatal if body cannot be unmarshalled.
func (req *Request) Unmarshal(v interface{}) {
	req.t.Helper()
	if req.Body == "" {
		req.t.Fatal(errors.New("golden file does not have body"))
		return
	}
	unmarshalBody(req.t, req.BodyType, req.Body, v)
}

// BindQuery binds request query parameters to v.
func (req *Request) BindQuery(tag string, v interface{}) {
	req.t.Helper()
	bindQuery(req.t, req.Query, tag, v)
}

// Bytes returns request body as byte slice.
func (req *Request) Bytes() []byte {
	return []byte(req.Body)
}
