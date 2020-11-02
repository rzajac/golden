package golden

import (
	"io"
	"regexp"
	"strings"
)

// Section represents golden file section.
type Section struct {
	// Identifier.
	id string

	// Lines as they appear in the golden file without
	// leading new line character.
	lines []string

	// Modifier.
	// Controls how to present the section. See parsing modifiers.
	mod string
}

// NewSection returns new Section instance or nil if the line is not a section.
// The NewSection will return non nil value only for lines with section
// identifiers.
func NewSection(line string) *Section {
	if isComment(line) {
		return NewComment(line)
	}

	name, mod, content := tokenize(line)
	if name == "" {
		return nil
	}

	sec := &Section{
		id:  line[:len(name)],
		mod: mod,
	}
	sec.lines = append(sec.lines, content)

	return sec
}

// NewComment returns special kind of section describing comments.
func NewComment(line string) *Section {
	sec := &Section{
		id:    SecComment,
		lines: []string{line[len(Delimiter):]},
		mod:   "",
	}
	return sec
}

// ID returns section identifier.
func (sec *Section) ID() string {
	return sec.id
}

// LineCnt returns number of lines in the section.
func (sec *Section) LineCnt() int {
	return len(sec.lines)
}

// Section returns section as a slice of bytes with additional new line
// character for the last line.
func (sec *Section) Section() []byte {
	out := sec.ID() + sec.mod + Delimiter + strings.Join(sec.lines, "\n") + "\n"
	return []byte(out)
}

// String implements fmt.Stringr interface and returns section as it was
// defined in the golden file. The last line doesn't have the new
// line character.
func (sec *Section) String() string {
	return strings.Join(sec.lines, modToSep(sec.mod))
}

// Add adds line to the section.
func (sec *Section) Add(line string) {
	_, _, line = tokenize(line)
	sec.lines = append(sec.lines, line)
}

// WriteTo implements io.WriteTo interface.
func (sec *Section) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(sec.Section())
	return int64(n), err
}

// secRX recognizes section start and modifiers.
var secRX = regexp.MustCompile(
	"^(" +
		SecRspCode + "|" +
		SecReqMethod + "|" +
		SecReqPath + "|" +
		SecReqQuery + "|" +
		SecHeader + "|" +
		SecBody + "|" +
		")(" +
		"." +
		")?" +
		Delimiter +
		"(.*)",
)

// tokenize checks if the line is matching the line which is the beginning
// of a section. Returns section name and modifier (if any) and the line
// content without section ID, modifier and delimiter.
func tokenize(line string) (string, string, string) {
	m := secRX.FindStringSubmatch(line)
	if m == nil {
		return "", "", line
	}
	return m[1], m[2], m[3]
}

// modToSep returns separator used to join lines based on modifier.
func modToSep(mod string) string {
	switch mod {
	case ModNone:
		return "\n"
	case ModMerge:
		return ""
	default:
		return "\n"
	}
}
