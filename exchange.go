package golden

import (
	"io"

	"gopkg.in/yaml.v3"
)

// Exchange represents HTTP request / response exchange.
type Exchange struct {
	// HTTP request.
	Request *Request `yaml:"request"`

	// HTTP response.
	Response *Response `yaml:"response"`

	// Test manager.
	t T
}

// NewExchange returns new instance of HTTP request / response Exchange.
func NewExchange(t T, data []byte) *Exchange {
	t.Helper()

	ex := &Exchange{}
	if err := yaml.Unmarshal(data, ex); err != nil {
		t.Fatal(err)
		return nil
	}
	ex.t = t

	if ex.Request != nil {
		ex.Request.t = t
		ex.Request.validate()
	}

	if ex.Response != nil {
		ex.Response.t = t
		ex.Response.validate()
	}

	return ex
}

// WriteTo writes golden file to w.
func (ex *Exchange) WriteTo(w io.Writer) (int64, error) {
	data, err := yaml.Marshal(ex)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(data)
	return int64(n), err
}
