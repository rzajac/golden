package _example

import (
	"testing"

	"github.com/rzajac/golden"
)

func Test_CompareData(t *testing.T) {
	// --- Given ---
	gld := golden.New(golden.Open(t, "../testdata/file.yaml", nil))

	// --- When ---
	data := []byte(`{
		"key1": "val1"
	}`)

	// --- Then ---
	gld.Assert(data)
}

type Data struct {
	Key1 string `json:"key1"`
}

func Test_Unmarshal(t *testing.T) {
	// --- Given ---
	gld := golden.New(golden.Open(t, "../testdata/file.yaml", nil))

	// --- When ---
	data := &Data{}
	gld.Unmarshall(data)

	// --- Then ---
	if data.Key1 != "val1" {
		t.Errorf("expected `%s` got `%s`", "val1", data.Key1)
	}
}
