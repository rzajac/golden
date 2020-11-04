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

	gld := &exchange{}
	if err := yaml.Unmarshal(data, gld); err != nil {
		t.Fatal(err)
		return nil
	}
	gld.t = t

	if gld.Request != nil {
		gld.Request.t = t
		gld.Request.Validate()
	}

	if gld.Response != nil {
		gld.Response.t = t
		gld.Response.Validate()
	}

	return gld
}

// WriteTo implements io.WriteTo interface for writing golden files.
func (gld *exchange) WriteTo(w io.Writer) (int64, error) {
	data, err := yaml.Marshal(gld)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(data)
	return int64(n), err
}
