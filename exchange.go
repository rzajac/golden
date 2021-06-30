package golden

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

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
func NewExchange(t T, r io.Reader) *Exchange {
	t.Helper()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
		return nil
	}

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

// Assert makes the request described in the golden file to host and asserts
// the response matches. It returns constructed request and received response
// in case further assertions need to be done.
func (ex *Exchange) Assert(host string) (*http.Request, *http.Response) {
	u := url.URL{
		Scheme:   ex.Request.Scheme,
		Host:     host,
		Path:     ex.Request.Path,
		RawQuery: ex.Request.Query,
	}

	req, err := http.NewRequest(
		ex.Request.Method,
		u.String(),
		strings.NewReader(ex.Request.Body),
	)
	if err != nil {
		ex.t.Fatal(err)
	}
	req.Header = lines2Headers(ex.t, ex.Request.Headers...)
	cli := &http.Client{}
	rsp, err := cli.Do(req)
	if err != nil {
		ex.t.Fatal(err)
	}
	ex.Response.Assert(rsp)
	return req, rsp
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
