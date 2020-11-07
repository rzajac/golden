package golden

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"gopkg.in/yaml.v3"
)

// File represents golden file with body and body type.
type File struct {
	BodyType string `yaml:"bodyType"`
	Body     string `yaml:"body"`
	t        T
}

// New returns golden File representation.
func New(t T, data []byte) *File {
	t.Helper()

	fil := &File{
		t: t,
	}
	if err := yaml.Unmarshal(data, fil); err != nil {
		t.Fatal(err)
		return nil
	}
	fil.t = t

	return fil
}

// Bytes returns body as byte slice.
func (fil *File) Bytes() []byte {
	return []byte(fil.Body)
}

// Assert asserts file body matches data. It chooses the bast way to
// compare two byte slices based on body type. For example when
// comparing JSON both byte slices don't have to be identical but
// they must represent the same data.
func (fil *File) Assert(data []byte) {
	fil.t.Helper()

	var equal bool
	switch fil.BodyType {
	case TypeJSON:
		equal = assertJSONEqual(fil.t, fil.Bytes(), data)
	case TypeText:
		equal = bytes.Equal(fil.Bytes(), data)
	default:
		equal = bytes.Equal(fil.Bytes(), data)
	}

	if !equal {
		fil.t.Fatalf(
			"expected request body to match want\n %s\ngot\n%s",
			fil.Body,
			string(data),
		)
		return
	}
}

// WriteTo writes golden file to w.
func (fil *File) WriteTo(w io.Writer) (int64, error) {
	data, err := yaml.Marshal(fil)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(data)
	return int64(n), err
}

// Unmarshall unmarshalls file body to v based on BodyType. Currently
// only JSON body type is supported.
func (fil *File) Unmarshall(v interface{}) {
	fil.t.Helper()
	if fil.Body != "" {
		if err := json.Unmarshal(fil.Bytes(), v); err != nil {
			fil.t.Fatal(err)
			return
		}
		return
	}
	fil.t.Fatal(errors.New("golden file empty body"))
}
