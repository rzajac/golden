package golden

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// File represents golden file with body and body type.
type File struct {
	Meta     map[string]interface{} `yaml:"meta,omitempty"`
	BodyType string                 `yaml:"bodyType"`
	Body     string                 `yaml:"body"`
	t        T
}

// New returns golden File representation.
func New(t T, r io.Reader) *File {
	t.Helper()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
		return nil
	}

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

// Unmarshal unmarshalls file body to v based on BodyType. When body type is
// set to text v can be pointer to sting or byte slice (with enough space to
// fit body). Calls Fatal if body cannot be unmarshalled.
func (fil *File) Unmarshal(v interface{}) {
	fil.t.Helper()
	if fil.Body == "" {
		fil.t.Fatal(errors.New("golden file empty body"))
		return
	}
	unmarshalBody(fil.t, fil.BodyType, fil.Body, v)
}
