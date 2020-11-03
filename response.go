package golden

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
)

// Response represents HTTP response backed by a golden file.
type Response struct {
	StatusCode int      `yaml:"statusCode"`
	Headers    []string `yaml:"headers"`
	Body       string   `yaml:"body"`
	JSONBody   bool     `yaml:"jsonBody"`

	headers http.Header // Request headers.
	t       T           // Test manager.
}

// Validate validates response loaded from golden file.
func (rsp *Response) Validate() {
	if rsp.StatusCode == 0 {
		rsp.t.Fatal(errors.New("HTTP response needs response code"))
	}

	if len(rsp.Headers) > 0 {
		rsp.headers = lines2Headers(rsp.t, rsp.Headers...)
	} else {
		rsp.headers = make(http.Header)
	}
}

// Assert asserts response matches the golden file.
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
	if rsp.JSONBody {
		equal = JSONBytesEqual(rsp.t, []byte(rsp.Body), body)
	} else {
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

// UnmarshallJSONBody unmarshalls request body to v. Calls Fatal if body
// section does not exist or json.Unmarshal returns error.
func (rsp *Response) UnmarshallJSONBody(v interface{}) {
	if rsp.Body != "" {
		if err := json.Unmarshal([]byte(rsp.Body), v); err != nil {
			rsp.t.Fatal(err)
		}
		return
	}
	rsp.t.Fatal(errors.New("golden file does not have body"))
}
