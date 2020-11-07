package golden

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
)

// Response represents golden file for HTTP response.
type Response struct {
	StatusCode int      `yaml:"statusCode"`
	Headers    []string `yaml:"headers"`
	BodyType   string   `yaml:"bodyType"`
	Body       string   `yaml:"body"`

	headers http.Header // Request headers.
	t       T           // Test manager.
}

// validate validates response loaded from golden file.
func (rsp *Response) validate() {
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
// To compare response bodies a method best suited for body type is used.
// For example when comparing JSON bodies both byte slices don't have to be
// identical but they must represent the same data.
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
	case TypeJSON:
		equal = assertJSONEqual(rsp.t, rsp.Bytes(), body)
	case TypeText:
		equal = bytes.Equal(rsp.Bytes(), body)
	default:
		equal = bytes.Equal(rsp.Bytes(), body)
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

// Unmarshall unmarshalls response body to v based on BodyType. Calls Fatal
// if body cannot be unmarshalled. Currently only JSON body type is supported.
func (rsp *Response) Unmarshall(v interface{}) {
	rsp.t.Helper()
	if rsp.Body != "" {
		if err := json.Unmarshal(rsp.Bytes(), v); err != nil {
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
