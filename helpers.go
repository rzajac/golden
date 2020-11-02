package golden

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"strings"
	"text/template"
)

// Open opens golden file.
func Open(t T, pth string) io.Reader {
	content, err := ioutil.ReadFile(pth)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(content)
}

// OpenTpl opens golden file template and renders it with data.
func OpenTpl(t T, pth string, data interface{}) io.Reader {
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

	return buf
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

// body2Lines returns http.Request body as lines.
func body2Lines(t T, req *http.Request) []string {
	buf := &bytes.Buffer{}
	tee := io.TeeReader(req.Body, buf)
	data, err := ioutil.ReadAll(tee)
	if err != nil {
		t.Fatal(err)
	}
	req.Body.Close()
	req.Body = ioutil.NopCloser(buf)

	lns := strings.Split(string(data), "\n")
	for i := range lns {
		lns[i] = strings.TrimRight(lns[i], "\r")
	}
	return lns
}
