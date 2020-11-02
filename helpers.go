package golden

import (
	"bytes"
	"io"
	"io/ioutil"
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
