package golden

import (
	"io"

	"gopkg.in/yaml.v3"
)

// T is a subset of testing.TB interface.
type T interface {
	// Fatal is equivalent to Log followed by FailNow.
	Fatal(args ...interface{})

	// Fatalf is equivalent to Logf followed by FailNow.
	Fatalf(format string, args ...interface{})

	// Helper marks the calling function as a test helper function.
	// When printing file and line information, that function will be skipped.
	// Helper may be called simultaneously from multiple goroutines.
	Helper()
}

// golden represents HTTP request / response golden file.
type golden struct {
	Request  *Request  `yaml:"request"`
	Response *Response `yaml:"response"`
	t        T         // Test manager.
}

// RequestResponse creates instance representing
// HTTP request / response golden file.
func RequestResponse(t T, data []byte) *golden {
	t.Helper()

	gld := &golden{}
	if err := yaml.Unmarshal(data, gld); err != nil {
		t.Fatal(err)
		return nil
	}
	gld.t = t

	if gld.Request != nil {
		gld.Request.t = t
		gld.Request.Validate()
	}

	if gld.Response != nil {
		gld.Response.t = t
		gld.Response.Validate()
	}

	return gld
}

// WriteTo implements io.WriteTo interface for writing golden files.
func (gld *golden) WriteTo(w io.Writer) (int64, error) {
	data, err := yaml.Marshal(gld)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(data)
	return int64(n), err
}
