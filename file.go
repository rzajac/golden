package golden

import (
	"encoding/json"
	"errors"
	"io"

	"gopkg.in/yaml.v3"
)

// file represents golden file with specific payload type.
type file struct {
	Payload     string `yaml:"payload"`
	PayloadType string `yaml:"payloadType"`
	t           T
}

// File unmarshalls YAML formatted data and returns new instance of file.
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

// Bytes returns payload as byte slice.
func (fil *file) Bytes() []byte {
	return []byte(fil.Payload)
}

// Assert asserts file payload matches data. It chooses the bast way to
// compare payloads based on payload data. For example when comparing JSON
// the data represented by JSON must match not the exact JSON string.
func (fil *file) Assert(data []byte) {
	fil.t.Helper()

	var equal bool
	switch fil.PayloadType {
	case PayloadJSON:
		equal = AssertJSONEqual(fil.t, []byte(fil.Payload), data)
	default:
		equal = fil.Payload == string(data)
	}

	if !equal {
		fil.t.Fatalf(
			"expected request body to match want\n %s\ngot\n%s",
			fil.Payload,
			string(data),
		)
		return
	}
}

// WriteTo writes golden file to w.
func (fil *file) WriteTo(w io.Writer) (int64, error) {
	data, err := yaml.Marshal(fil)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(data)
	return int64(n), err
}

// Unmarshall unmarshalls file payload to v based on PayloadType. Currently
// only JSON payload is supported.
func (fil *file) Unmarshall(v interface{}) {
	fil.t.Helper()
	if fil.Payload != "" {
		if err := json.Unmarshal([]byte(fil.Payload), v); err != nil {
			fil.t.Fatal(err)
			return
		}
		return
	}
	fil.t.Fatal(errors.New("golden file does not have payload"))
}
