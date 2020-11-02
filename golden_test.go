package golden

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Golden_basic(t *testing.T) {
	// --- When ---
	gld := NewGolden(t, Open(t, "testdata/request_basic.gold"))

	// --- Then ---
	assert.Exactly(t, 6, gld.SectionCount())

	sec := gld.Section(SecComment)
	require.NotNil(t, sec)
	assert.Exactly(t, []string{" Comment line 0.", " Comment line 1.", ""}, sec.lines)
	assert.Exactly(t, ModNone, sec.mod)

	sec = gld.Section(SecReqMethod)
	require.NotNil(t, sec)
	assert.Exactly(t, []string{"POST"}, sec.lines)
	assert.Exactly(t, ModNone, sec.mod)

	sec = gld.Section(SecReqPath)
	require.NotNil(t, sec)
	assert.Exactly(t, []string{"/some/path"}, sec.lines)
	assert.Exactly(t, ModNone, sec.mod)

	sec = gld.Section(SecReqQuery)
	require.NotNil(t, sec)
	assert.Exactly(t, []string{"key0=val0&key1=val1"}, sec.lines)
	assert.Exactly(t, ModNone, sec.mod)

	sec = gld.Section(SecHeader)
	require.NotNil(t, sec)
	assert.Exactly(t, []string{"Authorization: Bearer token"}, sec.lines)
	assert.Exactly(t, ModNone, sec.mod)

	sec = gld.Section(SecBody)
	require.NotNil(t, sec)
	assert.Exactly(t, []string{`{"key2": "val2"}`}, sec.lines)
	assert.Exactly(t, ModNone, sec.mod)
}

func Test_Golden_full_lines_modifier(t *testing.T) {
	// --- When ---
	gld := NewGolden(t, Open(t, "testdata/request_full.gold"))

	// --- Then ---
	assert.Exactly(t, 6, gld.SectionCount())

	sec := gld.Section(SecComment)
	require.NotNil(t, sec)
	exp := []string{" Comment line 0.", " Comment line 1.", ""}
	assert.Exactly(t, exp, sec.lines)
	assert.Exactly(t, ModNone, sec.mod)

	sec = gld.Section(SecReqMethod)
	require.NotNil(t, sec)
	assert.Exactly(t, []string{"POST"}, sec.lines)
	assert.Exactly(t, ModNone, sec.mod)

	sec = gld.Section(SecReqPath)
	require.NotNil(t, sec)
	assert.Exactly(t, []string{"/some/path"}, sec.lines)
	assert.Exactly(t, ModNone, sec.mod)

	sec = gld.Section(SecReqQuery)
	require.NotNil(t, sec)
	assert.Exactly(t, []string{"", "key0=val0", "&key1=val1", ""}, sec.lines)
	assert.Exactly(t, ModMerge, sec.mod)

	sec = gld.Section(SecHeader)
	require.NotNil(t, sec)
	exp = []string{
		"Authorization: Bearer token", "Content-Type: application/json",
		"",
	}
	assert.Exactly(t, exp, sec.lines)
	assert.Exactly(t, ModNone, sec.mod)

	sec = gld.Section(SecBody)
	require.NotNil(t, sec)
	assert.Exactly(t, []string{"", "{", `    "key2": "val2"`, "}"}, sec.lines)
	assert.Exactly(t, ModNone, sec.mod)
}

func Test_Golden_multi_line_body(t *testing.T) {
	// --- When ---
	gld := NewGolden(t, Open(t, "testdata/body_multi_line_json.gold"))

	// --- Then ---
	assert.Exactly(t, 2, gld.SectionCount())

	sec := gld.Section(SecBody)
	require.NotNil(t, sec)
	assert.Exactly(t, "\n{\n    \"key2\": \"val2\"\n}", sec.String())
	assert.Exactly(t, ModNone, sec.mod)

	exp := "Body::\n{\n    \"key2\": \"val2\"\n}\n"
	assert.Exactly(t, exp, string(sec.Section()))
}

func Test_Golden_Template(t *testing.T) {
	// --- Given ---
	data := map[string]interface{}{
		"val1": 1,
		"val2": "val2",
	}
	rdr := OpenTpl(t, "testdata/request_template.tpl.gold", data)

	// --- When ---
	gld := NewGolden(t, rdr)

	// --- Then ---
	sec := gld.Section(SecReqQuery)
	assert.NotNil(t, sec)
	assert.Exactly(t, "key0=val0&key1=1", sec.String())

	sec = gld.Section(SecBody)
	assert.NotNil(t, sec)
	assert.Exactly(t, `{"key2": "val2"}`, sec.String())
}

func Test_Golden_WriteTo(t *testing.T) {
	// --- Given ---
	gld := NewGolden(t, Open(t, "testdata/request_full.gold"))
	dst := &bytes.Buffer{}

	// --- When ---
	n, err := gld.WriteTo(dst)

	// --- Then ---
	assert.NoError(t, err)
	assert.Exactly(t, int64(215), n)

	exp, err := ioutil.ReadFile("testdata/request_full.gold")
	assert.NoError(t, err)
	assert.Exactly(t, string(exp), dst.String())
}
