package golden

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Section_NewSection_id(t *testing.T) {
	tt := []struct {
		testN string

		line string
		exp  string
	}{
		{"1", "Body::body", SecBody},
		{"2", "Body+:: body", SecBody},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			sec := NewSection(tc.line)

			// --- Then ---
			assert.Exactly(t, tc.exp, sec.ID(), "test %s", tc.testN)
		})
	}
}

func Test_Section_NewSection_lines(t *testing.T) {
	tt := []struct {
		testN string

		line string
		exp  []string
	}{
		{"1", "Body::body", []string{"body"}},
		{"2", "Body+:: body", []string{" body"}},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			sec := NewSection(tc.line)

			// --- Then ---
			assert.Exactly(t, tc.exp, sec.lines, "test %s", tc.testN)
		})
	}
}

func Test_Section_LineCnt(t *testing.T) {
	// --- Given ---
	sec := &Section{
		lines: []string{"A", "B"},
	}

	// --- Then ---
	assert.Exactly(t, 2, sec.LineCnt())
}

func Test_Section_NewSection_modifier(t *testing.T) {
	tt := []struct {
		testN string

		line string
		exp  string
	}{
		{"1", "Body::body", ModNone},
		{"2", "Body+:: body", ModMerge},
		{"3", "Body*:: body", "*"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			sec := NewSection(tc.line)

			// --- Then ---
			assert.Exactly(t, tc.exp, sec.mod, "test %s", tc.testN)
		})
	}
}

func Test_Section_NewSection_String(t *testing.T) {
	tt := []struct {
		testN string

		id    string
		lines []string
		mod   string
		exp   string
	}{
		{"1", SecBody, []string{}, ModNone, ""},
		{"2", SecBody, []string{"A", "B", "C"}, ModNone, "A\nB\nC"},
		{"3", SecBody, []string{"A", "B", "C"}, ModMerge, "ABC"},
		{"4", SecBody, []string{"A", "B", "C"}, "*", "A\nB\nC"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			sec := &Section{
				id:    tc.id,
				lines: tc.lines,
				mod:   tc.mod,
			}

			// --- When ---
			got := sec.String()

			// --- Then ---
			assert.Exactly(t, tc.exp, got, "test %s", tc.testN)
		})
	}
}

func Test_Section_NewSection_Section(t *testing.T) {
	tt := []struct {
		testN string

		id    string
		lines []string
		mod   string
		exp   []byte
	}{
		{"1", SecBody, []string{}, ModNone, []byte("Body::\n")},
		{"2", SecBody, []string{"A", "B", "C"}, ModNone, []byte("Body::A\nB\nC\n")},
		{"3", SecBody, []string{"A", "B", "C"}, ModMerge, []byte("Body+::A\nB\nC\n")},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			sec := &Section{
				id:    tc.id,
				lines: tc.lines,
				mod:   tc.mod,
			}

			// --- When ---
			got := sec.Section()

			// --- Then ---
			assert.Exactly(t, tc.exp, got, "test %s", tc.testN)
		})
	}
}

func Test_Section_NewSection_Write(t *testing.T) {
	tt := []struct {
		testN string

		id    string
		lines []string
		mod   string
		expN  int64
		exp   []byte
	}{
		{"1", SecBody, []string{}, ModNone, 7, []byte("Body::\n")},
		{"2", SecBody, []string{"A", "B", "C"}, ModNone, 12, []byte("Body::A\nB\nC\n")},
		{"3", SecBody, []string{"A", "B", "C"}, ModMerge, 13, []byte("Body+::A\nB\nC\n")},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			sec := &Section{
				id:    tc.id,
				lines: tc.lines,
				mod:   tc.mod,
			}
			buf := &bytes.Buffer{}

			// --- When ---
			n, err := sec.WriteTo(buf)

			// --- Then ---
			require.NoError(t, err, "test %s", tc.testN)
			assert.Exactly(t, tc.expN, n, "test %s", tc.testN)
			assert.Exactly(t, tc.exp, buf.Bytes(), "test %s", tc.testN)
		})
	}
}

func Test_Section_NewComment(t *testing.T) {
	tt := []struct {
		testN string

		line string
		exp  []string
	}{
		{"1", ":: Comment", []string{" Comment"}},
		{"2", "::  Comment", []string{"  Comment"}},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			sec := NewComment(tc.line)

			// --- Then ---
			assert.Exactly(t, SecComment, sec.id, "test %s", tc.testN)
			assert.Exactly(t, ModNone, sec.mod, "test %s", tc.testN)
			assert.Exactly(t, tc.exp, sec.lines, "test %s", tc.testN)
		})
	}
}

func Test_section_tokenize(t *testing.T) {
	tt := []struct {
		testN string

		line   string
		expSec string
		expMod string
		expCon string
	}{
		{"1", "RspCode::abc", SecRspCode, "", "abc"},
		{"2", "RspCode+:: abc", SecRspCode, ModMerge, " abc"},
		{"3", "not a section", "", "", "not a section"},
		{"4", "RspCode::ReqMethod::", SecRspCode, "", "ReqMethod::"},
		{"5", "RspCode*::abc", SecRspCode, "*", "abc"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			sec, mod, content := tokenize(tc.line)

			// --- Then ---
			assert.Exactly(t, tc.expSec, sec, "test %s", tc.testN)
			assert.Exactly(t, tc.expMod, mod, "test %s", tc.testN)
			assert.Exactly(t, tc.expCon, content, "test %s", tc.testN)
		})
	}
}
