package golden

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_File_New(t *testing.T) {
	// --- When ---
	gld := New(Open(t, "testdata/file.yaml", nil))

	// --- Then ---
	assert.Exactly(t, TypeJSON, gld.BodyType)
	assert.Exactly(t, `{ "key1": "val1" }`, gld.Body)
}

func Test_File_Bytes(t *testing.T) {
	// --- When ---
	gld := New(Open(t, "testdata/file.yaml", nil))

	// --- Then ---
	exp := []byte(`{ "key1": "val1" }`)
	assert.Exactly(t, exp, gld.Bytes())
}

func Test_File_WriteTo(t *testing.T) {
	// --- Given ---
	gld := New(Open(t, "testdata/file.yaml", nil))
	dst := &bytes.Buffer{}

	// --- When ---
	n, err := gld.WriteTo(dst)

	// --- Then ---
	assert.NoError(t, err)
	assert.Exactly(t, int64(42), n)

	got := New(t, dst)
	assert.Exactly(t, TypeJSON, got.BodyType)
	assert.Exactly(t, `{ "key1": "val1" }`, got.Body)
}

func Test_File_Unmarshall(t *testing.T) {
	// --- Given ---
	gld := New(Open(t, "testdata/file.yaml", nil))

	// --- When ---
	m := make(map[string]string, 1)
	gld.Unmarshall(&m)

	// --- Then ---
	require.Len(t, m, 1)
	require.Contains(t, m, "key1")
	assert.Exactly(t, "val1", m["key1"])
}
