package golden

import (
	"errors"
	"net/http"
)

// Response represents HTTP response backed by a golden file.
type Response struct {
	Code    int      `yaml:"code"`
	Headers []string `yaml:"headers"`
	Body    string   `yaml:"body"`

	headers http.Header // Request headers.
	t       T           // Test manager.
}

// Validate validates response loaded from golden file.
func (rsp *Response) Validate() {
	if rsp.Code == 0 {
		rsp.t.Fatal(errors.New("HTTP response needs response code"))
	}

	if len(rsp.Headers) > 0 {
		rsp.headers = lines2Headers(rsp.t, rsp.Headers...)
	} else {
		rsp.headers = make(http.Header)
	}
}
