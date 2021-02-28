package golden

import (
	"bytes"
	"testing"
	"time"

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

func Test_File_WriteTo_WithMetadata(t *testing.T) {
	// --- Given ---
	gld := New(Open(t, "testdata/file_metadata.yaml", nil))
	dst := &bytes.Buffer{}

	// --- When ---
	n, err := gld.WriteTo(dst)

	// --- Then ---
	assert.NoError(t, err)
	assert.Exactly(t, int64(144), n)

	got := New(t, dst)
	assert.Exactly(t, "Line1\nLine2\n{{ .data }}\n", got.Body)
	assert.Exactly(t, "val1", got.Meta["key1"])
	assert.Exactly(t, 123, got.Meta["key2"])
	assert.Exactly(t, 12.3, got.Meta["key3"])

	exp := time.Date(2021, 2, 28, 10, 24, 25, 123000000, time.UTC)
	assert.Exactly(t, exp, got.Meta["key4"].(time.Time))
}

func Test_File_Unmarshal_JSON(t *testing.T) {
	// --- Given ---
	gld := New(Open(t, "testdata/file.yaml", nil))

	// --- When ---
	m := make(map[string]string, 1)
	gld.Unmarshal(&m)

	// --- Then ---
	require.Len(t, m, 1)
	require.Contains(t, m, "key1")
	assert.Exactly(t, "val1", m["key1"])
}

func Test_File_WithMetadata(t *testing.T) {
	// --- Given ---
	data := map[string]interface{}{
		"data": "data line",
	}

	// --- When ---
	gld := New(Open(t, "testdata/file_metadata.yaml", data))

	// --- Then ---
	assert.NotNil(t, gld)
	assert.Exactly(t, "Line1\nLine2\ndata line\n", gld.Body)
	assert.Exactly(t, "val1", gld.Meta["key1"])
	assert.Exactly(t, 123, gld.Meta["key2"])
	assert.Exactly(t, 12.3, gld.Meta["key3"])

	exp := time.Date(2021, 2, 28, 10, 24, 25, 123000000, time.UTC)
	assert.Exactly(t, exp, gld.Meta["key4"].(time.Time))
}
