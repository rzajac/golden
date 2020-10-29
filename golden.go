package golden

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"net/url"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Golden file sections.
const (
	secVerb   = "Verb::"
	secPath   = "Path::"
	secQuery  = "Query::"
	secHeader = "Header::"
	secBody   = "Body::"
)

// fatalFn is a function signature for reporting fatal errors.
// It is only used in Golden structure testing that's why it's not exported.
type fatalFn func(args ...interface{})

// Opt represents function signature for Golden constructor.
type Opt func(gld *Golden)

// TplData is a Golden constructor function option setting template data.
func TplData(data interface{}) Opt {
	return func(gld *Golden) {
		gld.data = data
	}
}

// Golden represents golden file.
type Golden struct {
	Comments   []string      // Comments on the top of the file.
	Verb       string        // HTTP verb.
	VerbSet    bool          // Set to true when Verb was set in golden file.
	Path       string        // HTTP path.
	PathSet    bool          // Set to true when Path was set in golden file.
	Query      url.Values    // HTTP query.
	QuerySet   bool          // Set to true when Query was set in golden file.
	Headers    http.Header   // HTTP headers.
	HeadersSet bool          // Set to true when Headers were set in golden file.
	Body       *bytes.Buffer // HTTP body.
	data       interface{}   // Template data.
	fatal      fatalFn
	t          *testing.T
}

// New reads golden file at pth and creates new instance of Golden.
func New(t *testing.T, pth string, opts ...Opt) *Golden {
	t.Helper()

	gld := &Golden{
		fatal:   t.Fatal,
		Headers: make(http.Header, 0),
		t:       t,
	}

	for _, opt := range opts {
		opt(gld)
	}

	var src io.Reader
	if gld.data != nil {
		got, err := ioutil.ReadFile(pth)
		if err != nil {
			gld.fatal(err)
		}

		tpl, err := template.New("golden").Parse(string(got))
		if err != nil {
			gld.fatal(err)
		}

		buf := &bytes.Buffer{}
		if err := tpl.Execute(buf, gld.data); err != nil {
			gld.fatal(err)
		}
		src = buf
	} else {
		var err error
		src, err = os.Open(pth)
		if err != nil {
			gld.fatal(err)
		}
	}

	scn := bufio.NewScanner(src)
	for scn.Scan() {
		if err := gld.processLine(scn.Text()); err != nil {
			gld.fatal(err)
		}
	}

	if err := scn.Err(); err != nil {
		gld.fatal(err)
	}

	return gld
}

func (gld *Golden) processLine(lin string) error {
	var err error

	// If body is seen everything to the end
	// of the file is treated as body content.
	if gld.Body != nil {
		gld.Body.WriteString(lin + "\n")
		return nil
	}

	switch {
	case strings.HasPrefix(lin, "#"):
		// Set comments only if nothing else was set.
		if !(gld.VerbSet || gld.PathSet || gld.QuerySet || gld.HeadersSet || gld.Body != nil) {
			gld.Comments = append(gld.Comments, lin)
		}

	case strings.HasPrefix(lin, secVerb):
		gld.Verb = lin[len(secVerb):]
		gld.VerbSet = true

	case strings.HasPrefix(lin, secPath):
		gld.Path = lin[len(secPath):]
		gld.PathSet = true

	case strings.HasPrefix(lin, secQuery):
		if gld.Query, err = url.ParseQuery(lin[len(secQuery):]); err != nil {
			return err
		}
		gld.QuerySet = true

	case strings.HasPrefix(lin, secHeader):
		if err := addHeader(gld.Headers, lin[len(secHeader):]); err != nil {
			return err
		}
		gld.HeadersSet = true

	case strings.HasPrefix(lin, secBody):
		if gld.Body == nil {
			gld.Body = &bytes.Buffer{}
		}
		gld.Body.WriteString(lin[len(secBody):])
	}

	return nil
}

// Request returns HTTP request matching golden file. It panics on error.
func (gld *Golden) Request() *http.Request {
	body := gld.Body.Bytes()
	req := httptest.NewRequest(gld.Verb, gld.Path, bytes.NewReader(body))
	req.URL.RawQuery = gld.Query.Encode()
	for key, vv := range gld.Headers {
		for _, v := range vv {
			req.Header.Add(key, v)
		}
	}
	return req
}

// AssertRequest asserts request matches the golden file.
// Only the sections that were set are asserted.
func (gld *Golden) AssertRequest(req *http.Request) {
	if gld.VerbSet {
		assert.Exactly(gld.t, http.MethodPost, req.Method)
	}

	if gld.PathSet {
		assert.Exactly(gld.t, "/some/path", req.URL.Path)
	}

	if gld.QuerySet {
		assert.Exactly(gld.t, gld.Query.Encode(), req.URL.RawQuery)
	}

	if gld.HeadersSet {
		assert.Exactly(gld.t, gld.Headers, req.Header)
	}

	if gld.Body != nil {
		body, err := ioutil.ReadAll(req.Body)
		require.NoError(gld.t, err, "ReadAll")
		assert.Exactly(gld.t, gld.Body.Bytes(), body)
	}
}

func (gld *Golden) WriteTo(w io.Writer) (int64, error) {
	var total int64

	if len(gld.Comments) > 0 {
		str := strings.Join(gld.Comments, "\n") + "\n\n"
		n, err := w.Write([]byte(str))
		total += int64(n)
		if err != nil {
			return total, err
		}
	}

	if gld.VerbSet {
		str := secVerb + gld.Verb + "\n"
		n, err := w.Write([]byte(str))
		total += int64(n)
		if err != nil {
			return total, err
		}
	}

	if gld.PathSet {
		str := secPath + gld.Path + "\n"
		n, err := w.Write([]byte(str))
		total += int64(n)
		if err != nil {
			return total, err
		}
	}

	if gld.QuerySet {
		str := secQuery + gld.Query.Encode() + "\n"
		n, err := w.Write([]byte(str))
		total += int64(n)
		if err != nil {
			return total, err
		}
	}

	if gld.HeadersSet {
		for h, vv := range gld.Headers {
			for _, v := range vv {
				str := secHeader + h + ": " + v + "\n"
				n, err := w.Write([]byte(str))
				total += int64(n)
				if err != nil {
					return total, err
				}
			}
		}
	}

	if gld.Body != nil {
		tmp := append([]byte(secBody), gld.Body.Bytes()...)
		n, err := w.Write(tmp)
		total += int64(n)
		if err != nil {
			return total, err
		}
	}

	return total, nil
}

// addHeader adds header line to headers.
func addHeader(hs http.Header, lin string) error {
	hr := bufio.NewReader(strings.NewReader(lin + "\r\n\r\n"))
	tp := textproto.NewReader(hr)
	h, err := tp.ReadMIMEHeader()
	if err != nil {
		return err
	}
	for key, vv := range h {
		for _, v := range vv {
			hs.Add(key, v)
		}
	}
	return nil
}

func SaveRequest(pth string, req *http.Request, coms ...string) error {
	fil, err := os.Create(pth)
	if err != nil {
		return err
	}

	if len(coms) > 0 {
		str := strings.Join(coms, "\n") + "\n\n"
		if _, err = fil.WriteString(str); err != nil {
			return err
		}
	}

	str := secVerb + req.Method + "\n"
	if _, err = fil.WriteString(str); err != nil {
		return err
	}

	str = secPath + req.URL.Path + "\n"
	if _, err = fil.WriteString(str); err != nil {
		return err
	}

	str = secQuery + req.URL.RawQuery + "\n"
	if _, err = fil.WriteString(str); err != nil {
		return err
	}

	for h, vv := range req.Header {
		for _, v := range vv {
			str = secHeader + h + ": " + v + "\n"
			if _, err = fil.WriteString(str); err != nil {
				return err
			}
		}
	}

	var body []byte
	if req.Body != nil {
		var buf bytes.Buffer
		tee := io.TeeReader(req.Body, &buf)
		body, err = ioutil.ReadAll(tee)
		if err != nil {
			return err
		}
		_ = req.Body.Close()
		req.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
		if _, err := fil.Write([]byte(secBody)); err != nil {
			return err
		}
		if _, err := fil.Write(body); err != nil {
			return err
		}
	}

	return nil
}
