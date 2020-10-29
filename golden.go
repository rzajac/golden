package golden

import (
	"bufio"
	"bytes"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"strings"
	"testing"
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

// Golden represents golden file.
type Golden struct {
	Verb       string       // HTTP verb.
	VerbSet    bool         // Set to true when Verb was set in golden file.
	Path       string       // HTTP path.
	PathSet    bool         // Set to true when Path was set in golden file.
	Query      url.Values   // HTTP query.
	QuerySet   bool         // Set to true when Query was set in golden file.
	Headers    http.Header  // HTTP headers.
	HeadersSet bool         // Set to true when Headers were set in golden file.
	Body       bytes.Buffer // HTTP body.
	BodySet    bool         // Set to true when Body was set in golden file.
	fatal      fatalFn
}

// New reads golden file at pth and creates new instance of Golden.
func New(t *testing.T, pth string, opts ...Opt) *Golden {
	t.Helper()

	gld := &Golden{
		fatal:   t.Fatal,
		Headers: make(http.Header, 0),
	}

	for _, opt := range opts {
		opt(gld)
	}

	fil, err := os.Open(pth)
	if err != nil {
		gld.fatal(err)
	}

	scn := bufio.NewScanner(fil)
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

	if gld.BodySet {
		gld.Body.WriteString(lin + "\n")
		return nil
	}

	switch {
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
		gld.Body.WriteString(lin[len(secBody):])
		gld.BodySet = true

	case strings.HasPrefix(lin, "#"):
		// Do nothing.
	}

	return nil
}

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
