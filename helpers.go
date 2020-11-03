package golden

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"net/url"
	"reflect"
	"strings"
	"text/template"

	"github.com/gorilla/schema"
)

// Open opens golden file.
func Open(t T, pth string) []byte {
	data, err := ioutil.ReadFile(pth)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// Template opens golden file template and renders it with data.
func Template(t T, pth string, data interface{}) []byte {
	content, err := ioutil.ReadFile(pth)
	if err != nil {
		t.Fatal(err)
	}

	tpl, err := template.New("golden").Parse(string(content))
	if err != nil {
		t.Fatal(err)
	}

	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		t.Fatal(err)
	}

	return buf.Bytes()
}

// Headers2Lines creates lines representing http.Header.
func Headers2Lines(t T, hs http.Header) []string {
	buf := &bytes.Buffer{}
	if err := hs.Write(buf); err != nil {
		t.Fatal(err)
	}
	lns := strings.Split(buf.String(), "\r\n")
	if len(lns) > 0 {
		lns = lns[:len(lns)-1]
	}
	return lns
}

// JSONBytesEqual compares the JSON in two byte slices.
func JSONBytesEqual(t T, a, b []byte) bool {
	var ja, jb interface{}
	if err := json.Unmarshal(a, &ja); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &jb); err != nil {
		t.Fatal(err)
	}
	return reflect.DeepEqual(jb, ja)
}

// Map is a helper type for constructing template values.
type Map map[string]interface{}

// Add adds key and val to map and returns map for chaining.
func (m Map) Add(key string, val interface{}) Map {
	m[key] = val
	return m
}

// lines2Headers creates http.Header from header lines.
func lines2Headers(t T, lines ...string) http.Header {
	rdr := strings.NewReader(strings.Join(lines, "\r\n") + "\r\n\r\n")
	tp := textproto.NewReader(bufio.NewReader(rdr))
	hs, err := tp.ReadMIMEHeader()
	if err != nil {
		t.Fatal(err)
	}
	return http.Header(hs)
}

// readBody reads all from rc and returns read data and io.ReadCloser
// with the same data so it can be used to "reset" body of
// a http.Request or http.Response.
func readBody(t T, rc io.ReadCloser) ([]byte, io.ReadCloser) {
	buf := &bytes.Buffer{}
	tee := io.TeeReader(rc, buf)
	data, err := ioutil.ReadAll(tee)
	if err != nil {
		t.Fatal(err)
	}
	rc.Close()
	lns := strings.Split(string(data), "\n")
	for i := range lns {
		lns[i] = strings.TrimRight(lns[i], "\r")
	}
	return []byte(strings.Join(lns, "\n")), ioutil.NopCloser(buf)
}

// bindQuery binds HTTP query to v.
func bindQuery(t T, query, tag string, v interface{}) {
	vs, err := url.ParseQuery(query)
	if err != nil {
		t.Fatal(err)
	}
	dec := schema.NewDecoder()
	dec.SetAliasTag(tag)
	if err := dec.Decode(v, vs); err != nil {
		t.Fatal(err)
	}
}
