package golden

import (
	"encoding/json"
	"errors"
	"io"

	"gopkg.in/yaml.v3"
)

// file represents golden file for file.
type file struct {
	PayloadType string `yaml:"payloadType"`
	Payload     string `yaml:"payload"`
	t           T
}

// File creates instance representing file.
func File(t T, data []byte) *file {
	t.Helper()

	fil := &file{
		t: t,
	}
	if err := yaml.Unmarshal(data, fil); err != nil {
		t.Fatal(err)
		return nil
	}
	fil.t = t

	return fil
}

// Assert asserts file payload matches got.
func (fil *file) Assert(got []byte) {
	fil.t.Helper()

	var equal bool
	switch fil.PayloadType {
	case PayloadJSON:
		equal = JSONBytesEqual(fil.t, []byte(fil.Payload), got)
	default:
		equal = fil.Payload == string(got)
	}

	if !equal {
		fil.t.Fatalf(
			"expected request body to match want\n %s\ngot\n%s",
			fil.Payload,
			string(got),
		)
		return
	}
}

// WriteTo implements io.WriteTo interface for writing golden files.
func (fil *file) WriteTo(w io.Writer) (int64, error) {
	data, err := yaml.Marshal(fil)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(data)
	return int64(n), err
}

// Unmarshall unmarshalls file payload to v based on PayloadType. Calls Fatal
// if payload cannot be unmarshalled.
func (fil *file) Unmarshall(v interface{}) {
	fil.t.Helper()
	if fil.Payload != "" {
		if err := json.Unmarshal([]byte(fil.Payload), v); err != nil {
			fil.t.Fatal(err)
		}
		return
	}
	fil.t.Fatal(errors.New("golden file does not have payload"))
}
