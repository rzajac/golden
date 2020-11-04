package golden

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_File_File(t *testing.T) {
	// --- When ---
	gld := File(t, Open(t, "testdata/file.yaml"))

	// --- Then ---
	assert.Exactly(t, PayloadJSON, gld.PayloadType)
	assert.Exactly(t, `{ "key1": "val1" }`, gld.Payload)
}

func Test_File_WriteTo(t *testing.T) {
	// --- Given ---
	gld := File(t, Open(t, "testdata/file.yaml"))
	dst := &bytes.Buffer{}

	// --- When ---
	n, err := gld.WriteTo(dst)

	// --- Then ---
	assert.NoError(t, err)
	assert.Exactly(t, int64(48), n)

	got := File(t, dst.Bytes())
	assert.Exactly(t, PayloadJSON, got.PayloadType)
	assert.Exactly(t, `{ "key1": "val1" }`, got.Payload)
}

func Test_File_Unmarshall(t *testing.T) {
	// --- Given ---
	gld := File(t, Open(t, "testdata/file.yaml"))

	// --- When ---
	m := make(map[string]string, 1)
	gld.Unmarshall(&m)

	// --- Then ---
	require.Len(t, m, 1)
	require.Contains(t, m, "key1")
	assert.Exactly(t, "val1", m["key1"])
}
