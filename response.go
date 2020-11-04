package golden

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
)

// Response represents golden file for HTTP response.
type Response struct {
	StatusCode int      `yaml:"statusCode"`
	Headers    []string `yaml:"headers"`
	Body       string   `yaml:"body"`
	BodyType   string   `yaml:"bodyType"`

	headers http.Header // Request headers.
	t       T           // Test manager.
}

// Validate validates response loaded from golden file.
func (rsp *Response) Validate() {
	if rsp.StatusCode == 0 {
		rsp.t.Fatal(errors.New("HTTP response needs response code"))
		return
	}

	if len(rsp.Headers) > 0 {
		rsp.headers = lines2Headers(rsp.t, rsp.Headers...)
	} else {
		rsp.headers = make(http.Header)
	}
}

// Assert asserts response matches the golden file.
//
// All headers defined in the golden file must match exactly but passed
// response may have more headers then defined in the golden file.
//
// To compare response bodies the method best for defined body type is used.
// For example when comparing JSON the data represented by JSON must match not
// the exact JSON string.
func (rsp *Response) Assert(got *http.Response) {
	rsp.t.Helper()

	if rsp.StatusCode != got.StatusCode {
		rsp.t.Fatalf(
			"expected response status code %d got %d",
			rsp.StatusCode,
			got.StatusCode,
		)
		return
	}

	// Checks only headers set in golden file, got request may have more.
	for key, vv := range rsp.headers {
		g := got.Header.Values(key)
		if !reflect.DeepEqual(vv, g) {
			rsp.t.Fatalf(
				"expected response header %s values %v got %v",
				key,
				vv,
				g,
			)
			return
		}
	}

	body, rc := readBody(rsp.t, got.Body)
	defer func() { got.Body = rc }()

	var equal bool
	switch rsp.BodyType {
	case PayloadJSON:
		equal = AssertJSONEqual(rsp.t, []byte(rsp.Body), body)
	default:
		equal = rsp.Body == string(body)
	}

	if !equal {
		rsp.t.Fatalf(
			"expected response body to match want\n %s\ngot\n%s",
			rsp.Body,
			body,
		)
		return
	}
}

// UnmarshallBody unmarshalls response body to v based on BodyType. Calls Fatal
// if body cannot be unmarshalled. Currently only JSON payload is supported.
func (rsp *Response) UnmarshallBody(v interface{}) {
	rsp.t.Helper()
	if rsp.Body != "" {
		if err := json.Unmarshal([]byte(rsp.Body), v); err != nil {
			rsp.t.Fatal(err)
			return
		}
		return
	}
	rsp.t.Fatal(errors.New("golden file does not have body"))
}

// Bytes returns request body as byte slice.
func (rsp *Response) Bytes() []byte {
	return []byte(rsp.Body)
}
