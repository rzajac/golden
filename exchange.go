package golden

import (
	"io"

	"gopkg.in/yaml.v3"
)

// exchange represents HTTP request / response exchange.
type exchange struct {
	// HTTP request.
	Request *Request `yaml:"request"`

	// HTTP response.
	Response *Response `yaml:"response"`

	// Test manager.
	t T
}

// Exchange creates instance representing HTTP request / response exchange.
func Exchange(t T, data []byte) *exchange {
	t.Helper()

	ex := &exchange{}
	if err := yaml.Unmarshal(data, ex); err != nil {
		t.Fatal(err)
		return nil
	}
	ex.t = t

	if ex.Request != nil {
		ex.Request.t = t
		ex.Request.Validate()
	}

	if ex.Response != nil {
		ex.Response.t = t
		ex.Response.Validate()
	}

	return ex
}

// WriteTo implements io.WriteTo interface for writing golden files.
func (ex *exchange) WriteTo(w io.Writer) (int64, error) {
	data, err := yaml.Marshal(ex)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(data)
	return int64(n), err
}
